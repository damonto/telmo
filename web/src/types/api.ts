export type ApiResponse<T = unknown> = {
  data: T
}

export type EmptyObject = Record<string, never>
