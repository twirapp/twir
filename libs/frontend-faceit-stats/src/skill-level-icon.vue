<script lang="ts" setup>
import Icon1 from './assets/skill_level_1.svg?url'
import Icon10 from './assets/skill_level_10.svg?url'
import Icon2 from './assets/skill_level_2.svg?url'
import Icon3 from './assets/skill_level_3.svg?url'
import Icon4 from './assets/skill_level_4.svg?url'
import Icon5 from './assets/skill_level_5.svg?url'
import Icon6 from './assets/skill_level_6.svg?url'
import Icon7 from './assets/skill_level_7.svg?url'
import Icon8 from './assets/skill_level_8.svg?url'
import Icon9 from './assets/skill_level_9.svg?url'

defineProps<{ lvl: number }>()

const icons: Record<number, string> = {
	1: Icon1,
	2: Icon2,
	3: Icon3,
	4: Icon4,
	5: Icon5,
	6: Icon6,
	7: Icon7,
	8: Icon8,
	9: Icon9,
	10: Icon10,
}
</script>

<template>
	<img :src="icons[lvl]" class="skill-level-icon" />
</template>

<style>
.skill-level-icon {
	height: 30px;
  width: 30px;
}
</style>
