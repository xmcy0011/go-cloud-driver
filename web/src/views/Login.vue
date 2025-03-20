<template>
  <div class="login-container">
    <el-card class="login-card">
      <template #header>
        <h2 class="login-title">网盘系统登录</h2>
      </template>
      
      <div class="login-methods">
        <div class="qr-section">
          <div class="qr-wrapper" v-loading="qrLoading">
            <img v-if="qrCode" :src="qrCode" class="qr-code" alt="微信登录二维码"/>
            <div v-else-if="qrExpired" class="qr-expired">
              <p>二维码已过期</p>
              <el-button type="primary" @click="refreshQrCode">刷新二维码</el-button>
            </div>
          </div>
          <p class="qr-tip">请使用微信扫码登录</p>
        </div>

        <div class="divider">
          <span>或</span>
        </div>

        <el-form
          :model="loginForm"
          :rules="rules"
          label-width="0"
          class="account-section"
        >
          <el-form-item prop="username">
            <el-input
              v-model="loginForm.username"
              placeholder="用户名"
            >
            <template #prefix>
                <el-icon><User /></el-icon>
              </template>
            </el-input>
          </el-form-item>
          
          <el-form-item prop="password">
            <el-input
              v-model="loginForm.password"
              type="password"
              placeholder="密码"
              show-password
            >
              <template #prefix>
                <el-icon><Lock /></el-icon>
              </template>
            </el-input>
          </el-form-item>
          
          <el-form-item>
            <el-button 
              type="primary" 
              class="login-button"
              :loading="loading"
              @click="handleLogin"
            >
              账号密码登录
            </el-button>
            <div class="register-link">
              还没有账号？<el-link type="primary" @click="goToRegister">立即注册</el-link>
            </div>
          </el-form-item>
        </el-form>
      </div>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import { ref, reactive, onMounted, onUnmounted } from 'vue';
import { useRouter } from 'vue-router';
import { ElMessage } from 'element-plus';
import { User, Lock } from '@element-plus/icons-vue';

const router = useRouter();
const loading = ref(false);
const qrLoading = ref(false);
const qrCode = ref('');
const qrExpired = ref(false);
let checkLoginTimer = null;

const loginForm = reactive({
  username: '',
  password: ''
});

const rules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' }
  ]
};

// 获取微信登录二维码
const getQrCode = async () => {
  qrLoading.value = true;
  qrExpired.value = false;
  try {
    // TODO: 调用后端接口获取二维码
    // const response = await fetch('/api/login/wx/qrcode');
    // const data = await response.json();
    // qrCode.value = data.qrCode;
    
    // 模拟获取二维码
    qrCode.value = 'https://via.placeholder.com/200x200?text=WeChatQR';
    startCheckLogin();
  } catch (error) {
    ElMessage.error('获取二维码失败');
  } finally {
    qrLoading.value = false;
  }
};

// 定时检查扫码状态
const startCheckLogin = () => {
  checkLoginTimer = setInterval(async () => {
    try {
      // TODO: 调用后端接口检查登录状态
      // const response = await fetch('/api/login/wx/check');
      // const data = await response.json();
      // if (data.status === 'logged_in') {
      //   handleLoginSuccess(data.token);
      // } else if (data.status === 'expired') {
      //   handleQrExpired();
      // }
    } catch (error) {
      console.error('检查登录状态失败:', error);
    }
  }, 2000);

  // 模拟二维码过期
  setTimeout(() => {
    handleQrExpired();
  }, 60000);
};

const handleQrExpired = () => {
  qrExpired.value = true;
  qrCode.value = '';
  clearInterval(checkLoginTimer);
};

const refreshQrCode = () => {
  getQrCode();
};

const handleLoginSuccess = (token) => {
  localStorage.setItem('token', token);
  ElMessage.success('登录成功');
  router.push('/dashboard');
};

// 账号密码登录
const handleLogin = async () => {
  loading.value = true;
  try {
    // TODO: 实现登录逻辑
    handleLoginSuccess('dummy-token');
  } catch (error) {
    ElMessage.error('登录失败');
  } finally {
    loading.value = false;
  }
};

const goToRegister = () => {
  router.push('/register');
};

onMounted(() => {
  getQrCode();
});

onUnmounted(() => {
  if (checkLoginTimer) {
    clearInterval(checkLoginTimer);
  }
});
</script>

<style scoped lang="scss">
.login-container {
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: #f5f7fa;
  
  .login-card {
    width: 800px;
    
    .login-title {
      text-align: center;
      margin: 0;
      color: #303133;
    }

    .login-methods {
      display: flex;
      justify-content: space-between;
      align-items: flex-start;
      padding: 20px;

      .qr-section {
        flex: 1;
        text-align: center;
        padding: 0 40px;

        .qr-wrapper {
          width: 200px;
          height: 200px;
          margin: 0 auto 20px;
          display: flex;
          justify-content: center;
          align-items: center;
          border: 1px solid #EBEEF5;
          border-radius: 4px;

          .qr-code {
            width: 100%;
            height: 100%;
          }

          .qr-expired {
            text-align: center;
            color: #909399;

            p {
              margin-bottom: 10px;
            }
          }
        }

        .qr-tip {
          color: #909399;
          font-size: 14px;
        }
      }

      .divider {
        width: 1px;
        background-color: #DCDFE6;
        margin: 0 20px;
        height: 300px;
        position: relative;

        span {
          position: absolute;
          left: 50%;
          top: 50%;
          transform: translate(-50%, -50%);
          background-color: white;
          padding: 10px 0;
          color: #909399;
        }
      }

      .account-section {
        flex: 1;
        padding: 0 40px;

        .login-button {
          width: 100%;
        }

        .register-link {
          text-align: center;
          margin-top: 15px;
          font-size: 14px;
        }
      }
    }
  }
}
</style> 