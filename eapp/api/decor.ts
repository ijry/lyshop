import { get, post, put, del } from '@/utils/request'

export const getDecorVariants = () => get<any>('/decor/index/variants')

export const getDecorVariant = (key: string) => get<any>(`/decor/index?variant=${encodeURIComponent(key)}`)

export const updateDecorVariant = (key: string, components: any[]) =>
  put<any>(`/decor/index?variant=${encodeURIComponent(key)}`, { components })

export const publishDecorVariant = (key: string) =>
  post<any>(`/decor/index/publish?variant=${encodeURIComponent(key)}`)

export const copyDecorVariant = (payload: { source_key: string; new_key: string; name: string }) =>
  post<any>('/decor/index/copies', {
    from_variant_key: payload.source_key,
    new_variant_key: payload.new_key,
    new_variant_name: payload.name,
  })

export const renameDecorVariant = (key: string, name: string) =>
  put<any>(`/decor/index/variants/${encodeURIComponent(key)}`, { variant_name: name })

export const deleteDecorVariant = (key: string) =>
  del<any>(`/decor/index/variants/${encodeURIComponent(key)}`)
