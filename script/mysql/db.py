import pymysql
import logging
import os


def getLogger() -> logging.Logger:
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

def connectMySQL() -> pymysql.Connection:
    host = os.getenv('mysql_host')
    port = os.getenv('mysql_port')
    user = "root"
    password = os.getenv('mysql_pwd')
    db = "efast"

    conn = pymysql.connect(host=host, port=port, user=user, password=password,
                           database=db, charset='utf8mb4', cursorclass=pymysql.cursors.DictCursor)

    logger = getLogger()
    logger.info(msg='connect db, host=%s, port=%d'.format(host, port))

    with conn:
        logger.info(
            msg='connect db success, host=%s, port=%d'.format(host, port))

    return conn
