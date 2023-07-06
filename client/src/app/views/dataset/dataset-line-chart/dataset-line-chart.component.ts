import { Component, Input, OnInit, ViewChild } from '@angular/core';
import { ChartDataset, ChartOptions, Color } from 'chart.js';
import { BaseChartDirective } from 'ng2-charts';
import { Dataset, Ng2Dataset, Sensordata } from 'src/app/models/dataset';
import { DatasetsService } from 'src/app/services/datasets.service';
import { environment } from 'src/environments/environment';

@Component({
  selector: 'app-dataset-line-chart',
  templateUrl: './dataset-line-chart.component.html',
  styleUrls: ['./dataset-line-chart.component.scss'],
})
export class DatasetLineChartComponent implements OnInit {
  @Input() dataset: Dataset;
  @ViewChild(BaseChartDirective) private _chart;
  liveSensordata: Sensordata;
  enableLive: boolean = false;
  socket: any;
  data: Ng2Dataset;
  constructor(private datasetService: DatasetsService) {}
  lineChartData: ChartDataset[];
  lineChartLabels: string[];

  lineChartOptions = {
    responsive: true,
    animate: false,
    plugins: {
      legend: {
        position: 'right',
      },
      title: {
        display: true,
        text: 'Chart.js Line Chart',
      },
    },
  };
  lineChartLegend = true;
  lineChartPlugins = [];
  lineChartType = 'line';

  ngOnInit(): void {
    this.datasetService.LoadLineChartDatasetByReference(this.dataset.reference, 1000).subscribe((res) => {
      if (res) {
        this.lineChartLabels = res.labels;
        this.lineChartData = res.lineChartdataset;
        this.lineChartOptions.plugins.title.text = this.dataset.title;
      }
    });
  }
  runLive() {
    const socket = new WebSocket(`${environment.wsUrl}/api/live`);
    var id = this.dataset.id;
    socket.onopen = function (e) {
      console.log('WebSocket Opened');
      // socket.send(`dataset`);
    };
    this.clearChart();
    this.socket = socket;
    var self = this;
    socket.onmessage = function (e) {
      const data = JSON.parse(e.data) as Sensordata;
      // only showing current dataset
      if (self.dataset.id == data.dataset_id) {
        self.liveSensordata = data;
        self.updateChart(data);
      }
    };
  }
  updateChart(data: Sensordata) {
    this.lineChartData.forEach((dset, index) => {
      console.log('pushing data', data.data[index]);
      // adding new data points to chart
      dset.data.push(+data.data[index]);
      const now = new Date();
      this.lineChartLabels.push(now.toLocaleTimeString());
      if (dset.data.length > 10) {
        dset.data.shift();
        this.lineChartLabels.shift();
      }
    });
    this._chart.update();
  }
  clearChart() {
    this.lineChartLabels = [];
    this.lineChartData.forEach((dset, index) => {
      dset.data = [];
    });
  }
}
