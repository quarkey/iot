import { Component, Input, OnInit } from '@angular/core';

@Component({
  selector: 'app-status-dot',
  templateUrl: './status-dot.component.html',
  styleUrls: ['./status-dot.component.scss'],
})
export class StatusDotComponent implements OnInit {
  @Input() state: string;
  @Input() hideLabel?: boolean = false;
  constructor() {}

  ngOnInit(): void {}
}
