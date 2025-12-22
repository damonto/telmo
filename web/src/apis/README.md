# API Client

基于 VueUse `createFetch` 的简洁 API 封装。

## 使用方式

```typescript
import { useFetch } from '@/lib/fetch'

// GET 请求
const { data, error, isFetching } = await useFetch('modems').get().json()

// POST 请求
const { data } = await useFetch('users').post({ name: 'John' }).json()

// PUT 请求
const { data } = await useFetch('users/1').put({ name: 'Jane' }).json()

// DELETE 请求
const { data } = await useFetch('users/1').delete().json()
```

## API 定义

所有 API 统一在 `apis/` 目录定义：

```typescript
// apis/modem.ts
import { useFetch } from '@/lib/fetch'

export const useModemApi = () => {
  const getModems = () => useFetch<ModemListResponse>('modems').get().json()
  const getModemById = (id: string) => useFetch<ModemApiResponse>(`modems/${id}`).get().json()

  return { getModems, getModemById }
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
- X-Request-Time header

自动处理：

- 401 错误（清除 token）
- 开发环境日志

## 扩展

如需自定义，修改 `lib/fetch.ts` 中的 `beforeFetch`、`afterFetch`、`onFetchError`。
