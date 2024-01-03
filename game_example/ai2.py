import socket
from sys import argv

class AI2:
    def __init__(self):
        self.name = ""
        self.board = None  # 初始化棋盘

    def update(self, board,name):
        self.board = board  # 更新棋盘
        self.name=name

    def make_move(self):
        max_count = -1
        best_move = None
        for row in range(6):
            for col in range(7):
                if self.board[row][col] == 0:
                    count = self.count_connected(row, col)
                    if count > max_count:
                        max_count = count
                        best_move = [row, col]
        return best_move  # 返回能连接最多我方棋子的位置

    def count_connected(self, row, col):
        count = 0
        directions = [(0, 1), (1, 0), (1, 1), (1, -1)]  # 横、竖、主对角线、副对角线四个方向
        for dx, dy in directions:
            for d in [-1, 1]:  # 每个方向两边
                x, y = row + d * dx, col + d * dy
                while 0 <= x < 6 and 0 <= y < 7 and self.board[x][y] == self.name:
                    count += 1
                    x, y = x + d * dx, y + d * dy
        return count
    
    
if __name__ == "__main__":
    ai = AI2()
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