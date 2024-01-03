import random
import socket
from sys import argv

class AI1:
    def __init__(self):
        self.name = ""
        self.board = None  # 初始化棋盘

    def update(self, board,name):
        self.board = board  # 更新棋盘
        self.name=name

    def make_move(self):
        empty_positions = [(row, col) for row in range(6) for col in range(7) if self.board[row][col] == 0]
        if empty_positions:
            return random.choice(empty_positions)  # 随机选择一个没有棋子的位置
        else:
            return None  # 如果棋盘已满，返回None
        
        
if __name__ == "__main__":
    ai = AI1()
    # 1.创建socket
    tcp_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

    # 2. 链接服务器
    server_addr = ("127.0.0.1", 8901)
    tcp_socket.connect(server_addr)

    while True:
        # 接收消息
        data = tcp_socket.recv(1024)
        message = data.decode("gbk")

        # 检查消息是否为 "close"
        if message == "close":
            break
        else:
            # 分割消息
            values = message.split(",")
            name = int(values[0])
            board = [values[i:i+7] for i in range(1, 43, 7)]
            board = [[int(cell) for cell in row] for row in board]
            # 更新 AI
            ai.update(board, name)
            move=ai.make_move()
            # 将 move 转换为一个由 "," 连接的字符串
            move_str = ",".join(str(x) for x in move)
            # 发送 move_str
            tcp_socket.send(move_str.encode("gbk"))

    # 3. 关闭套接字
    tcp_socket.close()