import socket

def connect_to_server(host, port):
    try:
        client_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

        client_socket.connect((host, port))

        print("connected to server")

        message = "Hello from client!"
        client_socket.sendall(message.encode())
        
    except Exception as e:
        print("Error:", e)

    finally:
        client_socket.close()


host = "127.0.0.1"  
port = 8080        
connect_to_server(host, port)