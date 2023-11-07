# -*- coding: utf-8 -*-

import pymysql
import logging
import os
from env import *

def getLogger():
    # Set up logger
    logger = logging.getLogger(__name__)
    logger.setLevel(logging.INFO)
    formatter = logging.Formatter(
        "%(asctime)s - %(name)s - %(levelname)s - %(message)s")

    # Log to console
    handler = logging.StreamHandler()
    handler.setFormatter(formatter)
    logger.addHandler(handler)

    return logger

# 连接数据库，需要配置环境变量
# mysql_host: 数据地址
# mysql_port：端口
# mysql_pwd: 密码
def connectMySQL():
    host = mySQLHost
    port = mySQLport
    user = mySQLUser
    password = mySQLPassword
    db = mySQLDatabase

    if not host:
        host = "127.0.0.1"
    if not port:
        port = "3306"
    if not user:
        user = "root"

    conn = pymysql.connect(host=host, port=int(port), user=user, password=password,
                           database=db, charset='utf8mb4', cursorclass=pymysql.cursors.DictCursor)

    logger = getLogger()
    logger.info(msg='connect db, host=%s, port=%d'.format(host, port))

    with conn:
        logger.info(
            msg='connect db success, host=%s, port=%d'.format(host, port))

    return conn
