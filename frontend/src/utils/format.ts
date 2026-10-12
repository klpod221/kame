/**
 * Formats a byte number to human-readable size string (B, KB, MB, GB, TB).
 */
export function formatBytes(bytes: number): string {
  if (!bytes || bytes <= 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

/**
 * Formats an ISO date string to a localized date & time string.
 */
export function formatDate(dateStr: string): string {
  if (!dateStr) return 'N/A'
  try {
    const d = new Date(dateStr)
    return d.toLocaleDateString() + ' ' + d.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
  } catch {
    return dateStr
  }
}

/**
 * Formats a numeric rating to a fixed single-decimal string.
 */
export function formatRating(rating: number): string {
  if (!rating || rating <= 0) return '0.0'
  return rating.toFixed(1)
}
