<template>
  <div style="padding: 20px">
    <amis-renderer :schema="schema" />
    <router-view></router-view>
  </div>
</template>
<script>
import HomeApi from "@/api/modules/HomeApi";
import AMISRenderer from "@/components/AMISRenderer";

export default {
  components: {
    "amis-renderer": AMISRenderer,
  },
  computed: {},
  data: function () {
    return {
      // pins: [
      //   {
      //     "replication_factor_min": -1,
      //     "replication_factor_max": -1,
      //     "name": "",
      //     "mode": "recursive",
      //     "shard_size": 104857600,
      //     "data_shards": 4,
      //     "parity_shards": 2,
      //     "user_allocations": null,
      //     "expire_at": "0001-01-01T00:00:00Z",
      //     "metadata": null,
      //     "pin_update": null,
      //     "origins": [],
      //     "cid": "QmTgBV75rpqCE55aFyJxBMd1z8ennrr5yQH3SqbEgjUDME",
      //     "type": "pin",
      //     "allocations": [],
      //     "max_depth": -1,
      //     "reference": null,
      //     "timestamp": "2024-04-17T15:42:28Z"
      //   },
      //   {
      //     "replication_factor_min": -1,
      //     "replication_factor_max": -1,
      //     "name": "start.sh",
      //     "mode": "recursive",
      //     "shard_size": 104857600,
      //     "data_shards": 4,
      //     "parity_shards": 2,
      //     "user_allocations": null,
      //     "expire_at": "0001-01-01T00:00:00Z",
      //     "metadata": null,
      //     "pin_update": null,
      //     "origins": [],
      //     "cid": "QmVSsdAQ47EAtmURBetDkDYFK9TdD8p3pBEeX4fBENrq5j",
      //     "type": "pin",
      //     "allocations": [],
      //     "max_depth": -1,
      //     "reference": null,
      //     "timestamp": "2024-04-17T16:11:48Z"
      //   }
      // ],
      schema: {
        type: "page",
        title: "基于纠删码的IPFS存储优化",
        body: [
          {
            type: "html",
            html: `
            <div style='width: 100%; overflow-x: auto; text-align:center;'>
              <img src='cluster_architecture.png' style='height:600px; width: auto;' />
            </div>
            <br>
            <br>
          `
          },
          {
            type: "form",
            title: "使用纠删码上传文件",
            api: {
              method: "post",
              url: 'http://127.0.0.1:9094/add?data-shards=4&parity-shards=2&erasure=true&name=${file.name}&raw-leaves=true&shard=true&shard-size=${file.size/4}',
              requestAdaptor: function (api, context) {
                let url = new URL(api.url);
                let shardSize = url.searchParams.get('shard-size');
                shardSize = Math.round(parseFloat(shardSize)) + 256 * 1024;
                url.searchParams.set('shard-size', shardSize);
                api.url = url.toString();
                return {
                  ...api,
                };
              },
              adaptor: function (payload, response, api, context) {
                // 将payload字符串分割成多个对象
                let payloads = payload.split('\n').filter(str => str.trim() !== '').map(str => JSON.parse(str));
                // 返回处理后的数据
                console.log(payloads)
                return {
                  rows: payloads,
                };
              },
              dataType: "form-data",
              responseData: {
                rows: "${rows}"
              },
            },
            body: [
              {
                type: "input-file",
                name: "file",
                accept: "*",
                maxSize: 1048576000, // 例如最大 1000MB
                asBlob: true,
                drag: true
              },
              {
                type: "table",
                source: "${rows}",
                label: "数据",
                defaultParams: {
                  orderBy: "name",
                },
                columns: [
                  {
                    type: "text",
                    name: "name",
                    label: "名称",
                    sortable: true
                  },
                  {
                    type: "text",
                    name: "cid",
                    label: "CID",
                    searchable: true
                  },
                  {
                    type: "text",
                    name: "size",
                    label: "大小"
                  },
                  {
                    type: "text",
                    name: "allocations",
                    label: "存储节点"
                  }
                ],
                autoLoad: true
              }
            ],
          },
          {
            type: "form",
            title: "恢复并下载文件",
            controls: [
              {
                type: "text",
                name: "cid",
                required: true
              },
              {
                type: "button",
                label: "下载",
                actionType: "download",
                api: "get:http://localhost:8888/ecget?cid=${cid}",
              }
            ]
          },
          {
            type: "crud",
            name: "pinset",
            title: "Pin列表",
            api: {
              method: 'get',
              url: 'http://127.0.0.1:9094/pins',
              adaptor: function (payload, response, api, context) {
                if (typeof payload !== 'string') {
                  console.error('Payload is not a string:', payload);
                  return;
                }
                let payloads = payload.split('\n') // 假设每个 JSON 对象都在新的一行
                payloads = payloads.filter(object => object.trim() !== '').map(object => JSON.parse(object))
                for (let i = 0; i < payloads.length; i++) {
                  const newPeerMap = [];
                  for (const [key, value] of Object.entries(payloads[i].peer_map)) {
                    const newValue = { ...value, cluster_peer_id: key };
                    delete newValue.ipfs_peer_addresses;
                    delete newValue.ipfs_peer_id;
                    delete newValue.priority_pin;
                    delete newValue.timestamp;
                    delete newValue.error;
                    delete newValue.attempt_count;
                    if (newValue.status != "remote") {
                      newPeerMap.push(newValue);
                    }
                  }
                  payloads[i].peer_map = newPeerMap;
                  console.log(newPeerMap)
                }
                return {
                  items: payloads,
                };
              }
            }, columns: [
              {
                name: 'cid',
                label: 'CID',
                searchable: true
              },
              {
                name: 'name',
                label: '名称',
                sortable: true
              },
              {
                name: 'peer_map',
                label: '存储节点信息',
                type: "json",
                placeholder: "-",
                listItem: {
                  PeerName: "${peername}",
                  ClusterPeerID: "${cluster_peer_id}",
                  Status: "${status}",
                }
              },
              {
                name: 'created',
                label: '创建时间'
              }
            ],
            refresh: true, // 刷新按钮
            headerToolbar: [
              "reload"
            ],
          }
        ],
      },
    }
  },
  beforeCreate() { },
  async created() {
    // const response = await fetch('http://127.0.0.1:9094/allocations?filter=all')
    // console.log(response)
    // const text = await response.text()
    // const objects = text.split('\n') // 假设每个 JSON 对象都在新的一行
    // this.pins = objects.filter(object => object.trim() !== '').map(object => JSON.parse(object))
    // console.log(this.pins)
  },
  beforeMount() { },
  mounted() { },
  beforeUpdate() { },
  updated() { },
  activated() { }, //when keep-alive
  deactivated() { }, //when keep-alive
  beforeDestroy() { },
  destroyed() { },
  methods: {
    getDatList() {
      HomeApi.getDatList()
    },
  },
};
</script>
