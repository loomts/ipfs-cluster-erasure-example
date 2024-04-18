<template>
  <div class="menu-item">
    <template v-for="item in menu">
      <el-menu-item
        :index="item.name"
        :key="item.name"
        :class="{ 'submenu-title-noDropdown': !isNest }"
      >
        <!-- <i :class="[item.meta.icon || 'fa fa-navicon']"></i> -->
        <span slot="title">{{ item.meta.title }}</span>
      </el-menu-item>
      <el-menu-item
          v-for="child in item.children"
          :key="child.name"
          :index="child.name"
          class="child-menu-item"
        >
          <span slot="title">{{ child.meta.title }}</span>
        </el-menu-item>
    </template>
  </div>
</template>

<script>
export default {
  name: "SidebarItem",
  props: {
    menu: {
      type: Array,
    },
    isNest: {
      type: Boolean,
      default: false,
    },
  },
};
</script>

<style scoped>
.child-menu-item {
  margin-left: 20px;
}
</style>