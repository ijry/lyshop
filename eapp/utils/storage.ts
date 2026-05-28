export function getStorage<T = string>(key: string): T | '' {
  return (uni.getStorageSync(key) || '') as T | ''
}

export function setStorage(key: string, value: unknown) {
  uni.setStorageSync(key, value)
}

export function removeStorage(key: string) {
  uni.removeStorageSync(key)
}
