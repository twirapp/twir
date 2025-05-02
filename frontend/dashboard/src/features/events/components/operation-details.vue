<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { PlusIcon } from 'lucide-vue-next'

import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Switch } from '@/components/ui/switch'
import {
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import VariableInput from '@/components/variable-input.vue'
import { operationTypeSelectOptions, flatOperations } from '@/components/events/helpers'
import OperationFilter from './operation-filter.vue'

const props = defineProps<{
  operationIndex: number
  operation: any
  onAddFilter: (operationIndex: number) => void
  onRemoveFilter: (operationIndex: number, filterIndex: number) => void
}>()

const { t } = useI18n()
</script>

<template>
  <div class="space-y-4">
    <FormField
      v-slot="{ componentField }"
      :name="`operations.${operationIndex}.type`"
    >
      <FormItem>
        <FormLabel>{{ t('events.operationType') }}</FormLabel>
        <FormControl>
          <Select
            v-bind="componentField"
            :placeholder="t('events.selectOperationType')"
          >
            <SelectTrigger>
              <SelectValue :placeholder="t('events.selectOperationType')" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem
                v-for="option in operationTypeSelectOptions"
                :key="option.value"
                :value="option.value"
              >
                {{ option.label }}
              </SelectItem>
            </SelectContent>
          </Select>
        </FormControl>
        <FormMessage />
      </FormItem>
    </FormField>

    <div v-if="operation.type && flatOperations[operation.type]?.haveInput">
      <FormField
        v-slot="{ componentField }"
        :name="`operations.${operationIndex}.input`"
      >
        <FormItem>
          <FormLabel>{{ t('events.input') }}</FormLabel>
          <FormControl>
            <VariableInput v-bind="componentField" input-type="textarea" />
          </FormControl>
          <FormMessage />
        </FormItem>
      </FormField>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <FormField
        v-slot="{ componentField }"
        :name="`operations.${operationIndex}.delay`"
      >
        <FormItem>
          <FormLabel>{{ t('events.delay') }}</FormLabel>
          <FormControl>
            <Input 
              v-bind="componentField" 
              type="number" 
              min="0"
              :placeholder="t('events.delayPlaceholder')" 
            />
          </FormControl>
          <FormDescription>{{ t('events.delayDescription') }}</FormDescription>
          <FormMessage />
        </FormItem>
      </FormField>

      <FormField
        v-slot="{ componentField }"
        :name="`operations.${operationIndex}.repeat`"
      >
        <FormItem>
          <FormLabel>{{ t('events.repeat') }}</FormLabel>
          <FormControl>
            <Input 
              v-bind="componentField" 
              type="number" 
              min="0"
              :placeholder="t('events.repeatPlaceholder')" 
            />
          </FormControl>
          <FormDescription>{{ t('events.repeatDescription') }}</FormDescription>
          <FormMessage />
        </FormItem>
      </FormField>
    </div>

    <div v-if="operation.type && flatOperations[operation.type]?.additionalValues?.includes('useAnnounce')">
      <FormField
        v-slot="{ value, handleChange }"
        :name="`operations.${operationIndex}.useAnnounce`"
      >
        <FormItem class="flex flex-row items-center justify-between rounded-lg border p-3 shadow-sm">
          <div class="space-y-0.5">
            <FormLabel>{{ t('events.useAnnounce') }}</FormLabel>
            <FormDescription>
              {{ t('events.useAnnounceDescription') }}
            </FormDescription>
          </div>
          <FormControl>
            <Switch
              :checked="value"
              @update:checked="handleChange"
            />
          </FormControl>
        </FormItem>
      </FormField>
    </div>

    <div v-if="operation.type && flatOperations[operation.type]?.additionalValues?.includes('timeoutTime')">
      <FormField
        v-slot="{ componentField }"
        :name="`operations.${operationIndex}.timeoutTime`"
      >
        <FormItem>
          <FormLabel>{{ t('events.timeoutTime') }}</FormLabel>
          <FormControl>
            <Input 
              v-bind="componentField" 
              type="number" 
              min="0"
              :placeholder="t('events.timeoutTimePlaceholder')" 
            />
          </FormControl>
          <FormDescription>{{ t('events.timeoutTimeDescription') }}</FormDescription>
          <FormMessage />
        </FormItem>
      </FormField>
    </div>

    <div v-if="operation.type && flatOperations[operation.type]?.additionalValues?.includes('timeoutMessage')">
      <FormField
        v-slot="{ componentField }"
        :name="`operations.${operationIndex}.timeoutMessage`"
      >
        <FormItem>
          <FormLabel>{{ t('events.timeoutMessage') }}</FormLabel>
          <FormControl>
            <VariableInput v-bind="componentField" input-type="textarea" />
          </FormControl>
          <FormMessage />
        </FormItem>
      </FormField>
    </div>

    <div v-if="operation.type && flatOperations[operation.type]?.additionalValues?.includes('target')">
      <FormField
        v-slot="{ componentField }"
        :name="`operations.${operationIndex}.target`"
      >
        <FormItem>
          <FormLabel>{{ t('events.target') }}</FormLabel>
          <FormControl>
            <Input 
              v-bind="componentField" 
              :placeholder="t('events.targetPlaceholder')" 
            />
          </FormControl>
          <FormDescription>{{ t('events.targetDescription') }}</FormDescription>
          <FormMessage />
        </FormItem>
      </FormField>
    </div>

    <FormField
      v-slot="{ value, handleChange }"
      :name="`operations.${operationIndex}.enabled`"
    >
      <FormItem class="flex flex-row items-center justify-between rounded-lg border p-3 shadow-sm">
        <div class="space-y-0.5">
          <FormLabel>{{ t('events.operationEnabled') }}</FormLabel>
          <FormDescription>
            {{ t('events.operationEnabledDescription') }}
          </FormDescription>
        </div>
        <FormControl>
          <Switch
            :checked="value"
            @update:checked="handleChange"
          />
        </FormControl>
      </FormItem>
    </FormField>

    <!-- Filters section -->
    <div class="mt-6">
      <div class="flex justify-between items-center mb-4">
        <h3 class="text-lg font-medium">{{ t('events.filters') }}</h3>
        <Button variant="outline" size="sm" @click="onAddFilter(operationIndex)">
          <PlusIcon class="h-4 w-4 mr-2" />
          {{ t('events.addFilter') }}
        </Button>
      </div>

      <div v-if="operation.filters.length === 0" class="text-center py-4 border rounded-md">
        <p class="text-muted-foreground">{{ t('events.noFilters') }}</p>
      </div>

      <OperationFilter
        v-for="(filter, filterIndex) in operation.filters"
        :key="filterIndex"
        :operation-index="operationIndex"
        :filter-index="filterIndex"
        :on-remove="onRemoveFilter"
      />
    </div>
  </div>
</template>
