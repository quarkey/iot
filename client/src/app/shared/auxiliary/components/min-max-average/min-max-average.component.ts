import { Component, Input, OnInit } from '@angular/core';

@Component({
  selector: 'app-min-max-average',
  templateUrl: './min-max-average.component.html',
  styleUrls: ['./min-max-average.component.scss'],
})
export class MinMaxAverageComponent implements OnInit {
  @Input() report: any;
  @Input() title: string;
  constructor() {}

  ngOnInit(): void {}
}
