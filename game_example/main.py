import json
import socket
import random
from sys import argv
# 你的AI类需要有一个name属性，用于标识自己是哪个AI
# 你的AI类需要有一个make_move方法，用于在游戏中下棋
# 若下棋位置非法，视为跳过回合


class FourInARow:
    def __init__(self):
        self.board = [[0 for _ in range(7)] for _ in range(6)]  # 6行7列的棋盘
        self.winner= None  # 胜利者，如果没有胜利者则为None

    def is_game_over(self):
        # 检查每一行是否有连续的四个相同的非零元素
        for row in self.board:
            for i in range(4):
                if row[i] == row[i+1] == row[i+2] == row[i+3] != 0:
                    self.winner = row[i]
                    return True

        # 检查每一列是否有连续的四个相同的非零元素
        for col in range(7):
            for i in range(3):
                if self.board[i][col] == self.board[i+1][col] == self.board[i+2][col] == self.board[i+3][col] != 0:
                    self.winner = self.board[i][col]
                    return True

        # 检查主对角线是否有连续的四个相同的非零元素
        for row in range(3):
            for col in range(4):
                if self.board[row][col] == self.board[row+1][col+1] == self.board[row+2][col+2] == self.board[row+3][col+3] != 0:
                    self.winner = self.board[row][col]
                    return True

        # 检查副对角线是否有连续的四个相同的非零元素
        for row in range(3, 6):
            for col in range(4):
                if self.board[row][col] == self.board[row-1][col+1] == self.board[row-2][col+2] == self.board[row-3][col+3] != 0:
                    self.winner = self.board[row][col]
                    return True

        return False

    def make_move(self, move,name):
        row, col = move
        if 0 <= row < 6 and 0 <= col < 7 and self.board[row][col] == 0:  # 检查位置是否合法
            self.board[row][col] = name  # 更新棋盘
            return move
        else:
            return None  # 如果位置不合法，返回None

def write(game):
    winner = game.winner
    # result = {"moves": moves, "winner": winner}
    # print(json.dumps(result, indent=4))
    print(winner)

if __name__ == '__main__':

    # 创建 socket 对象
    server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

    # 绑定到指定地址和端口
    server_socket.bind(("127.0.0.1", 8901))

    # 开始监听，等待客户端连接
    server_socket.listen(2)

    # 接受两个客户端的连接
    client1, addr1 = server_socket.accept()
    client2, addr2 = server_socket.accept()

    # 随机决定先后顺序
    clients = [client1, client2]
    random.shuffle(clients)
    
    game = FourInARow()
    moves = []
    game_over = False

    # 轮流给客户端发送消息，等待回复
    while not game_over:
        for i, client in enumerate(clients):
            # 生成消息
            board_values = [str(cell) for row in game.board for cell in row]
            message = str(i+1) + "," + ",".join(board_values)

            # 发送消息
            client.send(message.encode("gbk"))

            # 等待回复
            data = client.recv(1024)
            
            # 分割回复并调用 make_move 函数
            row, col = map(int, data.decode("gbk").split(","))
            move = (row, col)
            game.make_move(move, i+1)

            # 将 [ai.name, move] 添加到 moves 数组中
            moves.append([i+1, move])

            # 检查游戏是否结束
            if game.is_game_over():
                game_over = True
                break

    # 关闭连接
    client1.send("close".encode("gbk"))
    client2.send("close".encode("gbk"))
    client1.close()
    client2.close()
    server_socket.close()
    write(game)