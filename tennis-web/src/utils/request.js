import axios from 'axios'

// Create an axios instance
const service = axios.create({
  baseURL: '/api', // Base URL for API requests. Adjust if needed (e.g. http://localhost:8080/api)
  timeout: 5000 // Request timeout
})

// Request interceptor
service.interceptors.request.use(
  config => {
    // Do something before request is sent
    return config
  },
  error => {
    // Do something with request error
    console.log(error) // for debug
    return Promise.reject(error)
  }
)

// Response interceptor
service.interceptors.response.use(
  response => {
    const res = response.data
    // You can add custom error handling here based on your API response structure
    return res
  },
  error => {
    console.log('err' + error) // for debug
    return Promise.reject(error)
  }
)

export default service
