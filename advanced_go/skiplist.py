import random
import math

class Node:
    def __init__(self, key, value=None, depth=1):
        self._key = key
        self._value = value
        # 一开始全部赋值为None
        self._next = [None for _ in range(depth)]
        self._depth = depth

    @property
    def key(self):
        return self._key

    @key.setter
    def key(self, key):
        self._key = key

    @property
    def value(self):
        return self._value

    @value.setter
    def value(self, value):
        self._value = value

    # 为第k个后向指针赋值
    def set_forward_pos(self, k, node):
        self._next[k] = node

    # 获取指定深度的指针指向的节点的key
    def query_key_by_depth(self, depth):
        # 后向指针指向的内容有可能为空，并且深度可能超界
        # 我们默认链表从小到大排列，所以当不存在的时候返回无穷大作为key
        return math.inf if depth > self._depth or self._next[depth] is None else self._next[depth].key

    # 获取指定深度的指针指向的节点
    def forward_by_depth(self, depth):
        return None if depth > self._depth else self._next[depth]

class SkipList:
    def __init__(self, max_depth, rate=0.5):
        # head的key设置成负无穷，tail的key设置成正无穷
        self.root = Node(-math.inf, depth=max_depth)
        self.tail = Node(math.inf)
        self.rate = rate
        self.max_depth = max_depth
        self.depth = 1
        # 把head节点的所有后向指针全部指向tail
        for i in range(self.max_depth):
            self.root.set_forward_pos(i, self.tail)

    def random_depth(self):
        depth = 1
        while True:
            rd = random.random()
            # 如果随机值小于p或者已经到达最大深度，就返回
            if rd < self.rate or depth == self.max_depth:
                return depth
            depth += 1

    def query(self, key):
        # 从头开始
        pnt = self.root
        # 遍历当下看的高度，高度只降不增
        for i in range(self.depth-1, -1, -1):
            # 如果看到比目标小的元素，则跳转
            while pnt.query_key_by_depth(i) < key:
                pnt = pnt.forward_by_depth(i)
        # 走到唯一可能出现的位置
        pnt = pnt.forward_by_depth(0)
        # 判断是否相等，如果相等则说明找到
        if pnt.key == key:
            return True, pnt.value
        else:
            return False, None

    def delete(self, key):
        # 记录下楼位置的数组
        heads = [None for _ in range(self.max_depth)]
        pnt = self.root
        for i in range(self.depth-1, -1, -1):
            while pnt.query_key_by_depth(i) < key:
                pnt = pnt.forward_by_depth(i)
            # 记录下楼位置
            heads[i] = pnt
        pnt = pnt.forward_by_depth(0)
        # 如果没找到，当然不存在删除
        if pnt.key == key:
            # 遍历所有下楼的位置
            for i in range(self.depth):
                # 由于是从低往高遍历，所以当看不到的时候，就说明已经超了，break
                if heads[i].forward_by_depth(i).key != key:
                    break
                # 将它看到的位置修改为删除节点同样楼层看到的位置
                heads[i].set_forward_pos(i, pnt.forward_by_depth(i))
            # 由于我们维护了skiplist当中的最高高度，所以要判断一下删除元素之后会不会出现高度降低的情况
            while self.depth > 1 and self.root.forward_by_depth(self.depth - 1) == self.tail:
                self.depth -= 1
        else:
            return False

    def insert(self, key, value):
        # 记录下楼的位置
        heads = [None for _ in range(self.max_depth)]
        pnt = self.root
        for i in range(self.depth-1, -1, -1):
            while pnt.query_key_by_depth(i) < key:
                pnt = pnt.forward_by_depth(i)
            heads[i] = pnt
        pnt = pnt.forward_by_depth(0)
        # 如果已经存在，直接修改
        if pnt.key == key:
            pnt.value = value
            return
        # 随机出楼层
        new_l = self.random_depth()
        # 如果楼层超过记录
        if new_l > self.depth:
            # 那么将头指针该高度指向它
            for i in range(self.depth, new_l):
                heads[i] = self.root
            # 更新高度
            self.depth = new_l

        # 创建节点
        new_node = Node(key, value, self.depth)
        for i in range(0, new_l):
            # x指向的位置定义成能看到x的位置指向的位置
            new_node.set_forward_pos(i, self.tail if heads[i] is None else heads[i].forward_by_depth(i))
            # 更新指向x的位置的指针
            if heads[i] is not None:
                heads[i].set_forward_pos(i, new_node)
