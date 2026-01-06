import axios from 'axios';
import { BASE_URL } from '../config/config';


// 创建 axios 实例，可以配置全局设置
const api = axios.create({
  baseURL: BASE_URL, // 设置基础 URL，根据实际情况调整
  timeout: 10000, // 设置请求超时时间
  headers: {
    'Content-Type': 'application/json',
  },
});

// 封装 GET 请求
export const get = async (url: string, params = {}) => {
  try {
    const response = await api.get(url, { params });
    return response.data; // 返回数据
  } catch (error) {
    throw error.response ? error.response.data : error.message;
  }
};

// 封装 POST 请求
export const post = async (url: string, data = {}) => {
  try {
    const response = await api.post(url, data);
    return response.data; // 返回数据
  } catch (error) {
    throw error.response ? error.response.data : error.message;
  }
};

// 你还可以根据需要封装 PUT、DELETE 请求
// 封装 PUT 请求
export const put = async (url: string, data = {}) => {
  try {
    const response = await api.put(url, data);
    return response.data;
  } catch (error) {
    throw error.response ? error.response.data : error.message;
  }
};

// 封装 DELETE 请求
export const remove = async (url: string, params = {}) => {
  try {
    const response = await api.delete(url, { params });
    return response.data;
  } catch (error) {
    throw error.response ? error.response.data : error.message;
  }
};
