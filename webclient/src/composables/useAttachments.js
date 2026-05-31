import { request, requestBlob } from '@/api/client'

const attachmentObjectUrls = new Set()

export function useAttachments() {
  async function hydrateMessageAttachments(message) {
    try {
      const attachments = await request(`/posts/${message.id}/attachments`)
      message.attachments = await Promise.all(
        attachments.map(async attachment => {
          const blob = await requestBlob(`/attachment/${encodeURIComponent(attachment.file_path)}`)
          const url = URL.createObjectURL(blob)
          attachmentObjectUrls.add(url)
          return { ...attachment, url, file_type: attachment.file_type || blob.type }
        })
      )
    } catch {
      message.attachments = []
    }
  }

  function isImageAttachment(attachment) {
    return attachment.file_type?.startsWith('image/')
  }

  function revokeAll() {
    attachmentObjectUrls.forEach(url => URL.revokeObjectURL(url))
    attachmentObjectUrls.clear()
  }

  return { hydrateMessageAttachments, isImageAttachment, revokeAll }
}
