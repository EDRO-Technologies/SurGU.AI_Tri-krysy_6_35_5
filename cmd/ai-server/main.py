from flask import Flask, request, jsonify
from pathlib import Path  
from rag import FastSafetyConsultantRAG  
import speech_recognition as sr

app = Flask(__name__)

rag_system = FastSafetyConsultantRAG(
    docs_directory="./documents",
    db_directory="./chroma_db",
    cache_directory="./cache",
    model_name="qwen2.5:7b",
    embedding_model="intfloat/multilingual-e5-small",
    enable_voice=False,  
    use_gpu=False
)

db_exists = (Path("./chroma_db") / "chroma.sqlite3").exists()
if not db_exists:
    print("Первый запуск. Обработка документов...")
    rag_system.process_documents()
else:
    print("Используем существующую БД.")
    rag_system.load_vectorstore()
    rag_system.setup_qa_chain()

@app.route('/ask', methods=['POST'])
def ask_question():
    data = request.get_json()  
    question = data.get('question')  

    if not question:
        return jsonify({"error": "No question provided"}), 400

    try:
        response = rag_system.ask(question)
        
        answer = response.get("answer")
        sources = response.get("sources", [])
        
        unique_files = set(source['file'] for source in sources)  

        files = [{"name": file_name} for file_name in unique_files]

        return jsonify({
            "answer": answer,
            "files": files
        })

    except Exception as e:
        return jsonify({"error": f"An error occurred: {str(e)}"}), 500



def speech_to_text(audio_data):
    recognizer = sr.Recognizer()
    audio = sr.AudioData(audio_data, 16000, 2)  

    try:
        text = recognizer.recognize_google(audio, language="ru-RU")
        return text
    except sr.UnknownValueError:
        return "Не удалось распознать речь"
    except sr.RequestError:
        return "Ошибка сервиса распознавания"


@app.route('/speech-to-text', methods=['POST'])
def speech_to_text_api():
    if 'audio' not in request.files:
        return jsonify({"error": "Отсутствует файл аудио"}), 400

    audio_file = request.files['audio']
    audio_data = audio_file.read()

    result = speech_to_text(audio_data)
    return jsonify({"text": result})


if __name__ == '__main__':
    app.run(debug=True, host='0.0.0.0', port=5000)
