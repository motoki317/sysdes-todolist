<template>
  <q-page padding class='row'>

    <todo-list></todo-list>

  </q-page>
</template>

<script lang='ts'>
import { defineComponent } from 'vue';
import TodoList from 'components/TodoList.vue';

export default defineComponent({
  name: 'PageIndex',

  components: { TodoList },

  setup() { return {} }
});
</script>
