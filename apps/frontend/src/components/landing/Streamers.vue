<script lang="ts" setup>
import Verified from '../../assets/icons/verified.svg?component';

const reviews = [
  { 'name': 'LWGerry', 'isVerified': false, 'message': 'Lorem ipsum dolor sit amet, consectetur adipiscing elit. ', 'subs': 245125 }, 
  { 'name': 'LWGerry', 'isVerified': false, 'message': 'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Proin eget nisi at mi elementum accumsan.  Integer sit amet viverra urna, nec posuere dui. Nullam bibendum lectus sed urna rhoncus vulputate. Aenean lectus lectus, rhoncus sit amet posuere sed, euismod in ligula. Sed rutrum ex sit amet nibh pellentesque, eget bibendum dui lacinia. Aliquam magna nibh, lacinia in luctus a, molestie quis magna.', 'subs': 245 }, 
  { 'name': 'LWGerry', 'isVerified': true, 'message': 'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Proin eget nisi at mi elementum accumsan. Proin posuere, lorem et euismod porta, nisi mauris sollicitudin enim, et placerat elit lectus in nulla. Integer sit amet viverra urna, nec posuere dui. Nullam bibendum lectus sed urna rhoncus vulputate. Aenean lectus lectus, rhoncus sit amet posuere sed, euismod in ligula. Sed rutrum ex sit amet nibh pellentesque, eget bibendum dui lacinia. Aliquam magna nibh, lacinia in luctus a, molestie quis magna.', 'subs': 245 }, 
  { 'name': 'LWGerry', 'isVerified': false, 'message': 'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Proin eget nisi at mi elementum accumsan. Proin posuere, lorem et euismod porta, nisi mauris sollicitudin enim, et placerat elit lectus in nulla. Integer sit amet viverra urna, nec posuere dui. Nullam bibendum lectus sed urna rhoncus vulputate. Aenean lectus lectus, rhoncus sit amet posuere sed, euismod in ligula. Sed rutrum ex sit amet nibh pellentesque, eget bibendum dui lacinia. Aliquam magna nibh, lacinia in luctus a, molestie quis magna.', 'subs': 2450 }, 
  { 'name': 'SadisNaMenya', 'isVerified': false, 'message': 'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Proin eget nisi at mi elementum accumsan. Proin posuere, lorem et euismod porta, nisi mauris sollicitudin enim, et placerat elit lectus in nulla. Integer sit amet viverra urna, nec posuere dui. Nullam bibendum lectus sed urna rhoncus vulputate. Aenean lectus lectus, rhoncus sit amet posuere sed, euismod in ligula. Sed rutrum ex sit amet nibh pellentesque, eget bibendum dui lacinia. Aliquam magna nibh, lacinia in luctus a, molestie quis magna.', 'subs': 245 }, 
  { 'name': 'LWGerry', 'isVerified': true, 'message': 'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Proin eget nisi at mi elementum accumsan.   Nullam bibendum lectus sed urna rhoncus vulputate. Aenean lectus lectus, rhoncus sit amet posuere sed, euismod in ligula. Sed rutrum ex sit amet nibh pellentesque, eget bibendum dui lacinia. Aliquam magna nibh, lacinia in luctus a, molestie quis magna.', 'subs': 24500 }, 
  { 'name': 'LWGerry', 'isVerified': false, 'message': 'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Proin eget nisi at mi elementum accumsan. Proin posuere, lorem et euismod porta, nisi mauris sollicitudin enim, et placerat elit lectus in nulla. Integer sit amet viverra urna, nec posuere dui. Nullam bibendum lectus sed urna rhoncus vulputate. Aenean lectus lectus, rhoncus sit amet posuere sed, euismod in ligula. Sed rutrum ex sit amet nibh pellentesque, eget bibendum dui lacinia. Aliquam magna nibh, lacinia in luctus a, molestie quis magna.', 'subs': 2450000 }, 
  { 'name': 'LWGerry', 'isVerified': false, 'message': 'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Proin eget nisi at mi elementum accumsan. Proin posuere, lorem et euismod porta, nisi mauris sollicitudin enim, et placerat elit lectus in nulla. Integer sit amet viverra urna, nec posuere dui. Nullam bibendum lectus sed urna rhoncus vulputate. Aenean lectus lectus, rhoncus sit amet posuere sed, euismod in ligula. Sed rutrum ex sit amet nibh pellentesque, eget bibendum dui lacinia. Aliquam magna nibh, lacinia in luctus a, molestie quis magna.', 'subs': 245 }, 
  { 'name': 'LWGerry', 'isVerified': false, 'message': 'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Proin eget nisi at mi elementum accumsan.  Integer sit amet viverra urna, nec posuere dui. lectus lectus, rhoncus sit amet posuere sed, euismod in ligula. Sed rutrum ex sit amet nibh pellentesque, eget bibendum dui lacinia. Aliquam magna nibh, lacinia in luctus a, molestie quis magna.', 'subs': 245 }, 
  { 'name': 'LWGerry', 'isVerified': false, 'message': 'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Proin eget nisi at mi elementum accumsan. Proin posuere, lorem et euismod porta, nisi mauris sollicitudin enim, et placerat elit lectus in nulla. Integer sit amet viverra urna, nec posuere dui. Nullam bibendum lectus sed urna rhoncus vulputate. Aenean lectus lectus, rhoncus sit amet posuere sed, euismod in ligula. Sed rutrum ex sit amet nibh pellentesque, eget bibendum dui lacinia. Aliquam magna nibh, lacinia in luctus a, molestie quis magna.', 'subs': 245 }] as Array<{'name': string, 'subs': number, 'isVerified':boolean, 'message': string}>;
</script>

<template>
  <div class="mb-32">
    <h2 class="font-bold mb-20 text-3xl text-center text-white">
      😎 Reviews from streamers
    </h2>
    
    <div class="lg:mx-40 md:mx-24 mx-2">
      <masonry-wall
        :items="reviews"
        :gap="30"
      >
        <template #default="{ item }">
          <div class="bg-[#181818] block drop-shadow-xl p-6 rounded text-white">
            <div class="flex flex-col">
              <div class="flex items-center mb-5 space-x-5">
                <img
                  src="https://i.picsum.photos/id/57/512/512.jpg?hmac=tup-ETUL0LHChWZTjAG31HXqECqaxgngmRKP87zOjDs"
                  alt="LOGO"
                  class="h-20 rounded-full w-20"
                >

                <h5 class="flex flex-col my-2 text-xl">
                  <div class="flex items-center space-x-2">
                    <span class="font-bold">{{ item.name }}</span>
                    <Verified />
                  </div>

                  <span class="text-gray-300 text-lg">{{ item.subs }} followers</span>
                </h5>
              </div>


              <p class="break-words my-2">
                {{ item.message }}
              </p>
            </div>
          </div>
        </template>
      </masonry-wall>
    </div>
  </div>
</template>