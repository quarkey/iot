import { Component, OnInit, Input } from '@angular/core';

@Component({
  selector: 'app-sensor-icon',
  templateUrl: './sensor-icon.component.html',
  styleUrls: ['./sensor-icon.component.scss'],
})
export class SensorIconComponent implements OnInit {
  @Input() category: string;
  constructor() {}

  ngOnInit(): void {}
}
