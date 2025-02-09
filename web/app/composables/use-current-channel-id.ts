export const useCurrentChannelId = () => useState<string | null>('currentChannelId', () => null)
export const setCurrentChannelId = (id: string | null) => useCurrentChannelId().value = id
