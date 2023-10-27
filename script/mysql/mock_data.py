
import time
from db import connectMySQL, getLogger

# if __name__ == "__main__":
conn = connectMySQL()
logger = getLogger()
logger.info(msg="hello")
