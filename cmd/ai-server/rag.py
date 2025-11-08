from pathlib import Path
from typing import List, Dict, Optional
from concurrent.futures import ThreadPoolExecutor, as_completed
import pickle

from langchain_community.document_loaders import PyPDFLoader, Docx2txtLoader, TextLoader
from langchain.text_splitter import RecursiveCharacterTextSplitter
from langchain_community.embeddings import HuggingFaceEmbeddings
from langchain_community.vectorstores import Chroma
from langchain_community.llms import Ollama
from langchain.chains import RetrievalQA
from langchain.prompts import PromptTemplate
from langchain.schema import Document

import speech_recognition as sr

from striprtf.striprtf import rtf_to_text

class FastSafetyConsultantRAG:
    def __init__(
        self,
        docs_directory: str = "./documents",
        db_directory: str = "./chroma_db",
        cache_directory: str = "./cache",
        model_name: str = "qwen2.5:7b",
        embedding_model: str = "intfloat/multilingual-e5-small",
        enable_voice: bool = True,
        use_gpu: bool = True
    ):
        self.docs_directory = Path(docs_directory)
        self.db_directory = Path(db_directory)
        self.cache_directory = Path(cache_directory)
        self.model_name = model_name
        self.embedding_model_name = embedding_model
        self.enable_voice = enable_voice
        self.use_gpu = use_gpu
        
        self.docs_directory.mkdir(exist_ok=True)
        self.db_directory.mkdir(exist_ok=True)
        self.cache_directory.mkdir(exist_ok=True)
    
        if self.enable_voice:
            self.recognizer = sr.Recognizer()
            self.recognizer.energy_threshold = 4000
            self.recognizer.dynamic_energy_threshold = True
            self.recognizer.pause_threshold = 0.8
        
        device = 'cuda' if use_gpu else 'cpu'
        self.embeddings = HuggingFaceEmbeddings(
            model_name=embedding_model,
            model_kwargs={'device': device},
            encode_kwargs={
                'normalize_embeddings': True,
                'batch_size': 32
            }
        )
        
        self.llm = Ollama(
            model=model_name,
            temperature=0.3,
            num_ctx=2048,
            num_predict=512,
            top_k=10,
            top_p=0.9
        )
        
        self.vectorstore = None
        self.qa_chain = None
    
    def _load_single_document(self, file_path: Path) -> List[Document]:
        try:
            suffix = file_path.suffix.lower()
            if suffix == '.pdf':
                loader = PyPDFLoader(str(file_path))
            elif suffix == '.docx':
                loader = Docx2txtLoader(str(file_path))
            elif suffix == '.txt':
                loader = TextLoader(str(file_path), encoding='utf-8')
            elif suffix == '.rtf':
                loader = RTFLoader(str(file_path))
            else:
                return []
            
            docs = loader.load()
            for doc in docs:
                doc.metadata['source_file'] = file_path.name
                doc.metadata['file_type'] = file_path.suffix
            
            return docs
        except Exception as e:
            return []
    
    def load_documents(self) -> List[Document]:
        cache_file = self.cache_directory / "documents_cache.pkl"
        if cache_file.exists():
            try:
                with open(cache_file, 'rb') as f:
                    documents = pickle.load(f)
                return documents
            except:
                pass
        
        supported_extensions = {'.pdf', '.docx', '.txt', '.rtf'}
        
        all_files = list(self.docs_directory.rglob('*'))
        
        files = [f for f in all_files if f.is_file() and f.suffix.lower() in supported_extensions]
        
        if not files:
            return []
        
        documents = []
        
        with ThreadPoolExecutor(max_workers=4) as executor:
            future_to_file = {
                executor.submit(self._load_single_document, f): f 
                for f in files
            }
            
            for future in as_completed(future_to_file):
                file_path = future_to_file[future]
                try:
                    docs = future.result()
                    if docs:
                        documents.extend(docs)
                except Exception as e:
                    pass
        
        with open(cache_file, 'wb') as f:
            pickle.dump(documents, f)
        
        return documents
    
    def split_documents(self, documents: List[Document]) -> List[Document]:
        text_splitter = RecursiveCharacterTextSplitter(
            chunk_size=800,
            chunk_overlap=100,
            length_function=len,
            separators=["\n\n", "\n", ". ", " ", ""]
        )
        
        return text_splitter.split_documents(documents)
    
    def create_vectorstore(self, chunks: List[Document]):
        batch_size = 100
        total_batches = (len(chunks) + batch_size - 1) // batch_size
        
        first_batch = chunks[:batch_size]
        self.vectorstore = Chroma.from_documents(
            documents=first_batch,
            embedding=self.embeddings,
            persist_directory=str(self.db_directory),
            collection_name="safety_docs"
        )
        
        for i in range(batch_size, len(chunks), batch_size):
            batch = chunks[i:i+batch_size]
            self.vectorstore.add_documents(batch)
    
    def load_vectorstore(self):
        self.vectorstore = Chroma(
            persist_directory=str(self.db_directory),
            embedding_function=self.embeddings,
            collection_name="safety_docs"
        )
    
    def setup_qa_chain(self):
        template = """Ты консультант по охране труда. Отвечай кратко и точно на основе документов.

Контекст: {context}

Вопрос: {question}

Ответ:"""

        prompt = PromptTemplate(
            template=template,
            input_variables=["context", "question"]
        )
        
        self.qa_chain = RetrievalQA.from_chain_type(
            llm=self.llm,
            chain_type="stuff",
            retriever=self.vectorstore.as_retriever(
                search_kwargs={"k": 3}
            ),
            return_source_documents=True,
            chain_type_kwargs={"prompt": prompt}
        )
    
    def ask(self, question: str) -> Dict:
        """Задать вопрос системе"""
        if not self.qa_chain:
            raise ValueError("Система не настроена")
        
        result = self.qa_chain.invoke({"query": question})
        
        response = {
            "answer": result["result"],
            "sources": []
        }
        
        for doc in result["source_documents"]:
            source_info = {
                "file": doc.metadata.get("source_file", "Неизвестно"),
                "page": doc.metadata.get("page", "N/A"),
                "content_preview": doc.page_content[:150] + "..."
            }
            response["sources"].append(source_info)
        
        return response
    
    def process_documents(self):
        documents = self.load_documents()
        
        if not documents:
            return
        
        chunks = self.split_documents(documents)
        self.create_vectorstore(chunks)
        self.setup_qa_chain()
