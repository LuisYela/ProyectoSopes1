const socket = io();

new Vue({
  el: '#ventana-cpu',
  created() {
    socket.on("statsMemory", (message) => {
      console.log(message)
      this.messagescpu.push(message)
      data.push({y:contador++,a:message.uso})
      console.log(data)
      grafico()
    })
  },
  data: {
    messagescpu: []
  },
  methods: {
    pedirStatscpu() {
      socket.emit("statscpu")
    }
  }
})
new Vue({
  el: '#ventana-ram',
  created() {
    socket.on("statsram", (message) => {
      console.log(message)
      this.messagesram.push(message)
      dataram1.push({y:contador,a:(message.disponible/10000)})
      dataram2.push({y:contador,a:(message.libre/10000)})
      dataram3.push({y:contador,a:(message.total/10000)})
      contador++
      console.log(data)
      graficoram()
    })
  },
  data: {
    messagesram: []
  },
  methods: {
    pedirStatsram() {
      socket.emit("statsram")
    }
  }
})
new Vue({
  el: '#ventana-stats',
  created() {
    socket.on("statsproc", (message) => {
      console.log(message)
      this.mproc=message
    }),
    socket.on("proclistado", (message) => {
      console.log(message)
      this.procs.push(message)
    })
  },
  data: {
    procesokill:"",
    mproc: {},
    procs: []
  },
  methods: {
    pedirStatsproc() {
      socket.emit("statsproc")
      this.cpustat = "";
    },
    borrarproc() {
      socket.emit("borrarproc", this.cpustat)
      this.cpustat = "";
    }
  }
})
var contador=0;
var data = [
]
function grafico() {
  var element, newElement, parent;
  element = document.getElementById("line-chart");
  parent = element.parentNode;
  newElement = document.createElement('div');
  newElement.id = "line-chart";
  parent.insertBefore(newElement, element);
  parent.removeChild(element);
  //newElement.innerHTML = "-new content here-";
  var datas = data,
  config = {
    data: datas,
    xkey: 'y',
    ykeys: ['a'],
    labels: ['Porcentaje de uso del CPU'],
    fillOpacity: 0.6,
    hideHover: 'auto',
    behaveLikeLine: true,
    resize: true,
    pointFillColors:['#ffffff'],
    pointStrokeColors: ['white'],
    lineColors:['red']
  };
  config.element = 'line-chart';
  Morris.Line(config)
}
var contador2=0;
var dataram1 = []
var dataram2 = []
var dataram3 = []
function graficoram() {
  var elementr1, elementr2, elementr3, newElementr1, newElementr2, newElementr3, parentr1, parentr2, parentr3;
  elementr1 = document.getElementById("line-chartram1");
  elementr2 = document.getElementById("line-chartram2");
  elementr3 = document.getElementById("line-chartram3");
  parentr1 = elementr1.parentNode;
  parentr2 = elementr2.parentNode;
  parentr3 = elementr3.parentNode;
  newElementr1 = document.createElement('div');
  newElementr2 = document.createElement('div');
  newElementr3 = document.createElement('div');
  newElementr1.id = "line-chartram1";
  newElementr2.id = "line-chartram2";
  newElementr3.id = "line-chartram3";
  parentr1.insertBefore(newElementr1, elementr1);
  parentr1.removeChild(elementr1);
  parentr2.insertBefore(newElementr2, elementr2);
  parentr2.removeChild(elementr2);
  parentr3.insertBefore(newElementr3, elementr3);
  parentr3.removeChild(elementr3);
  //newElement.innerHTML = "-new content here-";
  var datas1 = dataram1,
  config1 = {
    data: datas1,
    xkey: 'y',
    ykeys: ['a'],
    labels: ['RAM libre'],
    fillOpacity: 0.6,
    hideHover: 'auto',
    behaveLikeLine: true,
    resize: true,
    pointFillColors:['#ffffff'],
    pointStrokeColors: ['white'],
    lineColors:['red']
  };
  config1.element = 'line-chartram1';
  Morris.Line(config1)

  var datas2 = dataram2,
  config2 = {
    data: datas2,
    xkey: 'y',
    ykeys: ['a'],
    labels: ['RAM disponible'],
    fillOpacity: 0.6,
    hideHover: 'auto',
    behaveLikeLine: true,
    resize: true,
    pointFillColors:['#ffffff'],
    pointStrokeColors: ['white'],
    lineColors:['red']
  };
  config2.element = 'line-chartram2';
  Morris.Line(config2)

  var datas3 = dataram3,
  config3 = {
    data: datas3,
    xkey: 'y',
    ykeys: ['a'],
    labels: ['Porcentaje de uso del CPU'],
    fillOpacity: 0.6,
    hideHover: 'auto',
    behaveLikeLine: true,
    resize: true,
    pointFillColors:['#ffffff'],
    pointStrokeColors: ['white'],
    lineColors:['red']
  };
  config3.element = 'line-chartram3';
  Morris.Line(config3)
}