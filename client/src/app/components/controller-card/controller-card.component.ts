import { Component, Input, OnInit } from '@angular/core';
import { Controller } from 'src/app/models/controllers';
import { ControllersService } from 'src/app/services/controllers.service';

@Component({
  selector: 'app-controller-card',
  templateUrl: './controller-card.component.html',
  styleUrls: ['./controller-card.component.scss'],
})
export class ControllerCardComponent implements OnInit {
  @Input() citem: Controller;
  loading = false;
  constructor(private controllerService: ControllersService) {}

  ngOnInit(): void {}
  updateState() {
    this.loading = true;
    // http://localhost:6001/api/controller/4/switch/on
    if (this.citem.switch == 1) {
      this.controllerService.setContllerState(this.citem.id, 'off').subscribe((res) => {
        if (res) {
          this.loading = false;
          this.citem.switch = res.switch;
        }
      });
    }
    if (this.citem.switch == 0) {
      this.controllerService.setContllerState(this.citem.id, 'on').subscribe((res) => {
        if (res) {
          this.loading = false;
          this.citem.switch = res.switch;
        }
      });
    }
  }
}
