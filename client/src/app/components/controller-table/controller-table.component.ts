import { Component, Input, OnInit } from '@angular/core';
import { animate, state, style, transition, trigger } from '@angular/animations';
import { Controller } from 'src/app/models/controllers';
import { ControllersService } from 'src/app/services/controllers.service';

@Component({
  selector: 'app-controller-table',
  templateUrl: './controller-table.component.html',
  styleUrls: ['./controller-table.component.scss'],
  animations: [
    trigger('detailExpand', [
      state('collapsed', style({ height: '0px', minHeight: '0' })),
      state('expanded', style({ height: '*' })),
      transition('expanded <=> collapsed', animate('225ms cubic-bezier(0.4, 0.0, 0.2, 1)')),
    ]),
  ],
})
export class ControllerTableComponent implements OnInit {
  citems: Controller[];
  loading: boolean = true;
  alertLoading: boolean = true;

  constructor(private cs: ControllersService) {}
  dataSource: Controller[] = [];
  columnsToDisplay = ['type', 'category', 'title', 'description', 'switch', 'alert', 'active'];
  columnsToDisplayWithExpand = [...this.columnsToDisplay, 'expand'];
  expandedElement: Controller | null;
  ngOnInit(): void {
    this.cs.LoadControllersList().subscribe((res) => {
      if (res) {
        this.citems = res;
        this.dataSource = res;
        this.loading = false;
        this.alertLoading = false;
      }
    });
  }

  updateSwitchState(citem: Controller) {
    this.loading = true;
    // http://localhost:6001/api/controller/4/switch/on
    if (citem.switch == 1) {
      this.cs.setControllerSwitchState(citem.id, 'off').subscribe((res) => {
        if (res) {
          this.loading = false;
          citem.switch = res.switch;
        }
      });
    }
    if (citem.switch == 0) {
      this.cs.setControllerSwitchState(citem.id, 'on').subscribe((res) => {
        if (res) {
          this.loading = false;
          citem.switch = res.switch;
        }
      });
    }
  }
  updateAlertState(citem: Controller) {
    this.loading = true;
    // http://localhost:6001/api/controller/4/switch/on
    if (citem.alert == true) {
      this.cs.setControllerAlertState(citem.id, 'off').subscribe((res) => {
        if (res) {
          this.alertLoading = false;
          this.loading = false;
          citem.alert = res.alert;
        }
      });
    }
    if (citem.alert == false) {
      this.cs.setControllerAlertState(citem.id, 'on').subscribe((res) => {
        if (res) {
          this.alertLoading = false;
          this.loading = false;
          citem.alert = res.alert;
        }
      });
    }
  }
  updateState(citem: Controller) {
    this.loading = true;
    // http://localhost:6001/api/controller/4/switch/on
    if (citem.active == true) {
      this.cs.setControllerState(citem.id, 'off').subscribe((res) => {
        if (res) {
          citem.active = res.active;
          this.loading = false;
        }
      });
    }
    if (citem.active == false) {
      this.cs.setControllerState(citem.id, 'on').subscribe((res) => {
        if (res) {
          citem.active = res.active;
          this.loading = false;
        }
      });
    }
  }
}
