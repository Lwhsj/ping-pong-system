import { ref } from 'vue'

export function useVideoRecorder() {
  const mediaRecorder = ref(null)
  const chunks = ref([])
  const isRecording = ref(false)
  const error = ref(null)
  const stream = ref(null)

  const startRecording = async () => {
    try {
      if (!stream.value || !stream.value.active) {
         // Ensure previous tracks are stopped if any exist in a bad state
         if (stream.value) {
             stream.value.getTracks().forEach(t => t.stop())
         }
         stream.value = await navigator.mediaDevices.getUserMedia({ video: true })
      }
      
      const options = { mimeType: 'video/webm;codecs=vp8' }
      if (!MediaRecorder.isTypeSupported(options.mimeType)) {
        options.mimeType = 'video/webm'
        if (!MediaRecorder.isTypeSupported(options.mimeType)) {
            delete options.mimeType // Fallback to browser default
        }
      }

      mediaRecorder.value = new MediaRecorder(stream.value, options.mimeType ? options : undefined)
      
      chunks.value = []

      mediaRecorder.value.ondataavailable = (e) => {
        if (e.data && e.data.size > 0) {
          chunks.value.push(e.data)
        }
      }

      mediaRecorder.value.start(1000) // Slice every 1 second
      isRecording.value = true
      console.log('Video recording started')
    } catch (err) {
      console.error('Error starting video recording:', err)
      error.value = err.message
      isRecording.value = false
    }
  }

  const stopAndGetBlob = () => {
    return new Promise((resolve) => {
      if (!mediaRecorder.value || mediaRecorder.value.state === 'inactive') {
        resolve(chunks.value.length > 0 ? new Blob(chunks.value, { type: 'video/webm' }) : null)
        return
      }

      mediaRecorder.value.onstop = () => {
        const blob = new Blob(chunks.value, { type: 'video/webm' })
        resolve(blob)
      }
      
      mediaRecorder.value.stop()
      // Do NOT stop stream here, so we can restart quickly
      isRecording.value = false
    })
  }

  const stopRecording = () => {
    if (mediaRecorder.value && mediaRecorder.value.state !== 'inactive') {
      mediaRecorder.value.stop()
    }
    if (stream.value) {
      stream.value.getTracks().forEach(track => track.stop())
      stream.value = null
    }
    isRecording.value = false
    chunks.value = []
    console.log('Video recording stopped and stream released')
  }

  return {
    startRecording,
    stopRecording,
    stopAndGetBlob,
    isRecording,
    error
  }
}
