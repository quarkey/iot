import { Component, Input, OnInit } from '@angular/core';
import { DatasetsService } from 'src/app/services/datasets.service';

@Component({
  selector: 'app-display-value-box',
  templateUrl: './display-value-box.component.html',
  styleUrls: ['./display-value-box.component.scss'],
})
export class DisplayValueBoxComponent implements OnInit {
  @Input() value?: number = 0.0;
  @Input() label?: string = 'Label missing';
  @Input() icon?: string = 'thermostat';
  @Input() unit?: string;
  @Input() ref: string;

  loading = true;

  constructor(private datasetService: DatasetsService) {}

  ngOnInit(): void {
    this.datasetService.LoadDatasetByReference(this.ref, 1).subscribe((data) => {
      if (data) {
        console.log(data);
      }
    });
  }
}
