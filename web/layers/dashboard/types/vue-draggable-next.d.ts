declare module 'vue-draggable-next' {
  import { DefineComponent } from 'vue'
  
  export const VueDraggableNext: DefineComponent<{
    modelValue: any
    itemKey: string
    handle?: string
    tag?: string
    componentData?: any
    clone?: (original: any) => any
    move?: (evt: any, originalEvent: any) => boolean | undefined
    animation?: number
    group?: string | { name: string, pull?: string | boolean | Function, put?: string | boolean | Function | string[] }
    disabled?: boolean
    ghostClass?: string
    chosenClass?: string
    dragClass?: string
    dataTransfer?: any
    swapThreshold?: number
    invertSwap?: boolean
    direction?: string
    forceFallback?: boolean
    fallbackClass?: string
    fallbackOnBody?: boolean
    draggable?: string
    dragoverBubble?: boolean
    removeCloneOnHide?: boolean
    emptyInsertThreshold?: number
    lockAxis?: string
    lockOffset?: number | [number, number] | string | [string, string]
    scrollSensitivity?: number
    scrollSpeed?: number
    delay?: number
    delayOnTouchOnly?: boolean
    touchStartThreshold?: number
    multiDrag?: boolean
    multiDragKey?: string
    selectedClass?: string
    fallbackTolerance?: number
    setData?: (dataTransfer: any, dragEl: any) => void
    sort?: boolean
    filter?: string | Function
    preventOnFilter?: boolean
    revertClone?: boolean
    onStart?: (evt: any) => void
    onEnd?: (evt: any) => void
    onAdd?: (evt: any) => void
    onUpdate?: (evt: any) => void
    onSort?: (evt: any) => void
    onRemove?: (evt: any) => void
    onFilter?: (evt: any) => void
    onMove?: (evt: any, originalEvent: any) => boolean | undefined
    onClone?: (evt: any) => void
    onChange?: (evt: any) => void
    onChoose?: (evt: any) => void
    onUnchoose?: (evt: any) => void
  }>
}
