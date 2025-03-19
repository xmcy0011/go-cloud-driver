<template>
  <el-container class="dashboard-container">
    <el-header>
      <div class="header-content">
        <h2>网盘系统</h2>
        <el-button type="text" @click="handleLogout">退出登录</el-button>
      </div>
    </el-header>
    
    <el-container>
      <el-aside width="200px">
        <el-menu
          default-active="files"
          class="aside-menu"
        >
          <el-menu-item index="files">
            <el-icon><Document /></el-icon>
            <span>我的文件</span>
          </el-menu-item>
          <el-menu-item index="shared">
            <el-icon><Share /></el-icon>
            <span>共享文件</span>
          </el-menu-item>
        </el-menu>
      </el-aside>
      
      <el-main>
        <el-card>
          <template #header>
            <div class="file-header">
              <span>文件列表</span>
              <el-button type="primary">上传文件</el-button>
            </div>
          </template>
          
          <el-table :data="fileList">
            <el-table-column prop="name" label="文件名" />
            <el-table-column prop="size" label="大小" />
            <el-table-column prop="updateTime" label="修改时间" />
            <el-table-column label="操作">
              <template #default="scope">
                <el-button type="text">下载</el-button>
                <el-button type="text">分享</el-button>
                <el-button type="text" class="delete-btn">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { Document, Share } from '@element-plus/icons-vue';

const router = useRouter();

const fileList = ref([
  {
    name: '示例文档.docx',
    size: '2.5MB',
    updateTime: '2023-12-20 10:30'
  },
  {
    name: '项目计划.xlsx',
    size: '1.8MB',
    updateTime: '2023-12-19 15:45'
  }
]);

const handleLogout = () => {
  localStorage.removeItem('token');
  router.push('/login');
};
</script>

<style scoped lang="scss">
.dashboard-container {
  height: 100vh;
  
  .el-header {
    background-color: #409EFF;
    color: white;
    padding: 0 20px;
    
    .header-content {
      height: 60px;
      display: flex;
      justify-content: space-between;
      align-items: center;
    }
  }
  
  .aside-menu {
    height: 100%;
    border-right: solid 1px #e6e6e6;
  }
  
  .file-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
  
  .delete-btn {
    color: #F56C6C;
  }
}
</style> 