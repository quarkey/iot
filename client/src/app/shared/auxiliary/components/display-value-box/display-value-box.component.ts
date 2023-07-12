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
  @Input() title?: string = 'Datapoint value';
  @Input() unit?: string;
  @Input() time: string;
  errorMessage: string;
  showError = false;

  constructor() {}

  ngOnInit(): void {
    // cheking if incoming value is a number, if not set it to 0.0
    if (isNaN(this.value)) {
      this.errorMessage = 'Input value is not a number';
      this.showError = true;
    }
  }
}
