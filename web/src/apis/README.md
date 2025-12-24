# API Client

基于 VueUse `createFetch` 的简洁 API 封装。

## 使用方式

```typescript
import { useFetch } from '@/lib/fetch'

// GET 请求 - 注意 data.value?.data 才是实际数据
const { data } = await useFetch<ModemListResponse>('modems').get().json()
const modems = data.value?.data // 提取实际数据

// POST 请求
const { data } = await useFetch('users').post({ name: 'John' }).json()

// PUT 请求
const { data } = await useFetch('users/1').put({ name: 'Jane' }).json()

// DELETE 请求
const { data } = await useFetch('users/1').delete().json()
```

## 响应格式

所有 API 响应都使用统一的 `ApiResponse<T>` 包装格式：

```typescript
// 成功响应
{
  "data": T  // 实际数据
}

// 错误响应
{
  "code": 404,
  "message": "modem with ID xxx not found"
}
```

## API 定义

所有 API 统一在 `apis/` 目录定义：

```typescript
// apis/modem.ts
import { useFetch } from '@/lib/fetch'

export const useModemApi = () => {
  const getModems = () => useFetch<ModemListResponse>('modems').get().json()
  const getModem = (id: string) => useFetch<ModemDetailResponse>(`modems/${id}`).get().json()

  return { getModems, getModem }
}
```

## 配置

在 `.env.development` 中配置：

```bash
VITE_API_BASE_URL=http://10.10.10.101:9527/api/v1
```

## 拦截器

自动添加：

- Authorization header（如果有 token）
- Content-Type header

自动处理：

- 401 错误（清除 token）
- 非 2xx 状态码（弹出错误提示）
- 网络错误（弹出错误提示）
- 开发环境日志

## 错误处理

所有 API 错误都会通过统一的错误处理器处理，自动使用 AlertDialog 组件显示错误信息。

错误处理逻辑：

1. **非 2xx 状态码**：自动显示后端返回的错误消息
2. **401 Unauthorized**：清除 token 并提示用户重新登录
3. **404 Not Found**：显示资源未找到提示
4. **500+ Server Error**：显示服务器错误提示
5. **网络错误**：显示网络错误提示

错误对话框通过 `App.vue` 中的 `ErrorAlert` 组件全局渲染，开发者无需在每个 API 调用处手动处理错误。

## 扩展

如需自定义，修改 `lib/fetch.ts` 中的 `beforeFetch`、`afterFetch`、`onFetchError`。
