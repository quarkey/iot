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
      }
    });
  }

  updateState(citem: Controller) {
    this.loading = true;
    // http://localhost:6001/api/controller/4/switch/on
    if (citem.switch == 1) {
      this.cs.setContllerState(citem.id, 'off').subscribe((res) => {
        if (res) {
          this.loading = false;
          citem.switch = res.switch;
        }
      });
    }
    if (citem.switch == 0) {
      this.cs.setContllerState(citem.id, 'on').subscribe((res) => {
        if (res) {
          this.loading = false;
          citem.switch = res.switch;
        }
      });
    }
  }
}
