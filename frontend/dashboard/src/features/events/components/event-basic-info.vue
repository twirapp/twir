<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { computed } from 'vue'
import { useForm } from 'vee-validate'

import { useCommandsApi } from '@/api/commands/commands'
import { useKeywordsApi } from '@/api/keywords'
import { Input } from '@/components/ui/input'
import {
  Form,
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
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { Switch } from '@/components/ui/switch'
import TwitchRewardsSelector from '@/components/rewardsSelector.vue'
import { eventTypeSelectOptions } from '@/components/events/helpers'

const props = defineProps<{
  form: ReturnType<typeof useForm>
}>()

const { t } = useI18n()

// Fetch commands and keywords for selectors
const commandsApi = useCommandsApi()
const { data: commandsData } = commandsApi.useQueryCommands()
const commands = computed(() => commandsData.value?.commands || [])

const keywordsApi = useKeywordsApi()
const { data: keywordsData } = keywordsApi.useQueryKeywords()
const keywords = computed(() => keywordsData.value?.keywords || [])
</script>

<template>
  <Card>
    <CardHeader>
      <CardTitle>{{ t('events.basicInfo') }}</CardTitle>
      <CardDescription>{{ t('events.basicInfoDescription') }}</CardDescription>
    </CardHeader>
    <CardContent class="space-y-4">
      <FormField
        v-slot="{ componentField }"
        name="type"
      >
        <FormItem>
          <FormLabel>{{ t('events.type') }}</FormLabel>
          <FormControl>
            <Select
              v-bind="componentField"
              :placeholder="t('events.selectType')"
            >
              <SelectTrigger>
                <SelectValue :placeholder="t('events.selectType')" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem
                  v-for="option in eventTypeSelectOptions"
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

      <FormField
        v-slot="{ componentField }"
        name="description"
      >
        <FormItem>
          <FormLabel>{{ t('events.description') }}</FormLabel>
          <FormControl>
            <Input v-bind="componentField" :placeholder="t('events.descriptionPlaceholder')" />
          </FormControl>
          <FormMessage />
        </FormItem>
      </FormField>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <FormField
          v-slot="{ value, handleChange }"
          name="enabled"
        >
          <FormItem class="flex flex-row items-center justify-between rounded-lg border p-3 shadow-sm">
            <div class="space-y-0.5">
              <FormLabel>{{ t('events.enabled') }}</FormLabel>
              <FormDescription>
                {{ t('events.enabledDescription') }}
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

        <FormField
          v-slot="{ value, handleChange }"
          name="onlineOnly"
        >
          <FormItem class="flex flex-row items-center justify-between rounded-lg border p-3 shadow-sm">
            <div class="space-y-0.5">
              <FormLabel>{{ t('events.onlineOnly') }}</FormLabel>
              <FormDescription>
                {{ t('events.onlineOnlyDescription') }}
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

      <div v-if="form.values.type === 'REDEMPTION_CREATED'">
        <FormField
          v-slot="{ componentField }"
          name="rewardId"
        >
          <FormItem>
            <FormLabel>{{ t('events.reward') }}</FormLabel>
            <FormControl>
              <TwitchRewardsSelector v-bind="componentField" />
            </FormControl>
            <FormMessage />
          </FormItem>
        </FormField>
      </div>

      <div v-if="form.values.type === 'COMMAND_USED'">
        <FormField
          v-slot="{ componentField }"
          name="commandId"
        >
          <FormItem>
            <FormLabel>{{ t('events.command') }}</FormLabel>
            <FormControl>
              <Select
                v-bind="componentField"
                :placeholder="t('events.selectCommand')"
              >
                <SelectTrigger>
                  <SelectValue :placeholder="t('events.selectCommand')" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem
                    v-for="command in commands"
                    :key="command.id"
                    :value="command.id"
                  >
                    {{ command.name }}
                  </SelectItem>
                </SelectContent>
              </Select>
            </FormControl>
            <FormMessage />
          </FormItem>
        </FormField>
      </div>

      <div v-if="form.values.type === 'KEYWORD_USED'">
        <FormField
          v-slot="{ componentField }"
          name="keywordId"
        >
          <FormItem>
            <FormLabel>{{ t('events.keyword') }}</FormLabel>
            <FormControl>
              <Select
                v-bind="componentField"
                :placeholder="t('events.selectKeyword')"
              >
                <SelectTrigger>
                  <SelectValue :placeholder="t('events.selectKeyword')" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem
                    v-for="keyword in keywords"
                    :key="keyword.id"
                    :value="keyword.id"
                  >
                    {{ keyword.text }}
                  </SelectItem>
                </SelectContent>
              </Select>
            </FormControl>
            <FormMessage />
          </FormItem>
        </FormField>
      </div>
    </CardContent>
  </Card>
</template>
