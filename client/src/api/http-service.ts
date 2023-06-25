import axios from 'axios';

const axiosInstance = axios.create({
  baseURL: '',
  withCredentials: true,
  headers: {
    'Content-Type': 'application/json;charset=utf-8',
  },
});

export const httpService = {
  get: axiosInstance.get,
  post: axiosInstance.post,
};
