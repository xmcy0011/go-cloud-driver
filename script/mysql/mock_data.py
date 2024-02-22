
import time
from db import connectMySQL, getLogger

def AddClosure(cursor, parentId, objectId):
    # 2. 把查询出的行的后代改为要插入的节点 id
    # 1. 查询出后代是 B3 的所有行
    # 3. 加上节点本身，深度为1"
    sql = """
    insert into metadata_closure(ancestor, descendant, depth)
    select t.ancestor,%s,t.depth+1 from metadata_closure as t 
    where t.descendant = %s
    union all select %s,%s,0"""
    cursor.execute(sql, (objectId, parentId, objectId, objectId))

if __name__ == "__main__":
    conn = connectMySQL()
    logger = getLogger()

    # test
    # |- a
    # |- b
    with conn.cursor() as cursor:
        AddClosure(cursor, "test", "test")
        AddClosure(cursor, "test", "a")
        AddClosure(cursor, "test", "b")

    conn.commit()