import { Component, Input, OnInit } from '@angular/core';

@Component({
  selector: 'app-min-max-average',
  templateUrl: './min-max-average.component.html',
  styleUrls: ['./min-max-average.component.scss'],
})
export class MinMaxAverageComponent implements OnInit {
  @Input() report: HighLowAverage;
  @Input() title: string;
  constructor() {}

  ngOnInit(): void {}
}
interface HighLowAverage {
  description: string;
  date_from: string;
  date_to: string;
  high_date: string;
  low_date: string;
  high: number;
  low: number;
  average: number;
  datapoints: number;
}
