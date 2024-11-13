import json
import mysql.connector

# Substitua 'dados.json' pelo caminho do seu arquivo JSON
with open('logs.json', 'r') as f:
    data = json.load(f)

# Conectar ao banco de dados MySQL
conn = mysql.connector.connect(
    host="localhost",
    user="root",  # Substitua pelo seu usuário MySQL
    password="Lf200397@@",  # Substitua pela sua senha MySQL
    database="logs_db"
)

cursor = conn.cursor()

# Inserir dados na tabela
for item in data:
    cursor.execute("""
    INSERT INTO logs (state, latency_in, processing_time, latency_out, download_time, status_code, response_size, total_time)
    VALUES (%s, %s, %s, %s, %s, %s, %s, %s)
    """, (
        item['state'], item['latency_in'], item['processing_time'], item['latency_out'],
        item['download_time'], item['status_code'], item['response_size'], item['total_time']
    ))

# Confirmar as alterações e fechar a conexão
conn.commit()
cursor.close()
conn.close()

print("Dados importados com sucesso!")
