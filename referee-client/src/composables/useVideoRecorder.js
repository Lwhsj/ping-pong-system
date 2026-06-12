import { ref } from 'vue'

export function useVideoRecorder() {
  const mediaRecorder = ref(null)
  const chunks = ref([])
  const isRecording = ref(false)
  const error = ref(null)
  const stream = ref(null)

  const recorderOptions = () => {
    const options = { mimeType: 'video/webm;codecs=vp8' }
    if (!MediaRecorder.isTypeSupported(options.mimeType)) {
      options.mimeType = 'video/webm'
      if (!MediaRecorder.isTypeSupported(options.mimeType)) {
        return undefined
      }
    }
    return options
  }

  const startRecording = async () => {
    try {
      error.value = null

      if (mediaRecorder.value && mediaRecorder.value.state !== 'inactive') {
        isRecording.value = true
        return true
      }

      if (!stream.value || !stream.value.active) {
         // Ensure previous tracks are stopped if any exist in a bad state
         if (stream.value) {
             stream.value.getTracks().forEach(t => t.stop())
         }
         stream.value = await navigator.mediaDevices.getUserMedia({ video: true })
      }

      mediaRecorder.value = new MediaRecorder(stream.value, recorderOptions())
      
      chunks.value = []

      mediaRecorder.value.ondataavailable = (e) => {
        if (e.data && e.data.size > 0) {
          chunks.value.push(e.data)
        }
      }

      mediaRecorder.value.onerror = (event) => {
        console.error('Video recording error:', event.error || event)
        error.value = event.error?.message || '视频录制失败'
        isRecording.value = false
      }

      // Keep each rally as one complete media container. MediaRecorder chunks are
      // collected only to assemble the final Blob after stop(), not as video files.
      mediaRecorder.value.start()
      isRecording.value = true
      console.log('Video recording started')
      return true
    } catch (err) {
      console.error('Error starting video recording:', err)
      error.value = err.message
      isRecording.value = false
      return false
    }
  }

  const stopAndGetBlob = () => {
    return new Promise((resolve) => {
      const recorder = mediaRecorder.value
      if (!recorder) {
        resolve(null)
        return
      }

      const blobType = recorder.mimeType || 'video/webm'

      if (recorder.state === 'inactive') {
        const blob = chunks.value.length > 0 ? new Blob(chunks.value, { type: blobType }) : null
        resolve(blob && blob.size > 0 ? blob : null)
        return
      }

      recorder.onstop = () => {
        const blob = chunks.value.length > 0 ? new Blob(chunks.value, { type: blobType }) : null
        isRecording.value = false
        resolve(blob && blob.size > 0 ? blob : null)
      }
      
      recorder.stop()
      // Do NOT stop stream here, so we can restart quickly
      isRecording.value = false
    })
  }

  const stopRecording = () => {
    const recorder = mediaRecorder.value
    if (recorder && recorder.state !== 'inactive') {
      recorder.ondataavailable = null
      recorder.onstop = null
      recorder.stop()
    }
    mediaRecorder.value = null
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
