import { animate, state, style, transition, trigger } from '@angular/animations';
import { Component, Input, OnInit } from '@angular/core';
import { Dataset } from 'src/app/models/dataset';

@Component({
  selector: 'app-dataset-table',
  templateUrl: './dataset-table.component.html',
  styleUrls: ['./dataset-table.component.scss'],
  animations: [
    trigger('detailExpand', [
      state('collapsed', style({ height: '0px', minHeight: '0' })),
      state('expanded', style({ height: '*' })),
      transition('expanded <=> collapsed', animate('225ms cubic-bezier(0.4, 0.0, 0.2, 1)')),
    ]),
  ],
})
export class DatasetTableComponent implements OnInit {
  @Input() dataSource: Dataset[];
  columnsToDisplay = ['title', 'sensor_title', 'description', 'telemetry'];
  columnsToDisplayWithExpand = [...this.columnsToDisplay, 'expand'];
  expandedElement: Dataset | null;
  constructor() {}

  ngOnInit(): void {}
}
