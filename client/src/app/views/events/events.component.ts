import { Component, OnInit, ViewChild } from '@angular/core';
import { EventsList } from 'src/app/models/event';
import { EventsService } from 'src/app/services/events.service';
import { MatSort, Sort } from '@angular/material/sort';
import { LiveAnnouncer } from '@angular/cdk/a11y';
import { MatTableDataSource } from '@angular/material/table';

@Component({
  selector: 'app-events',
  templateUrl: './events.component.html',
  styleUrls: ['./events.component.scss'],
})
export class EventsComponent implements OnInit {
  constructor(private eventsService: EventsService, private _liveAnnouncer: LiveAnnouncer) {}
  displayedColumns: string[] = ['id', 'category', 'message', 'event_time'];
  dataSource: any;
  loading: boolean = true;
  @ViewChild(MatSort) sort: MatSort;
  ngOnInit(): void {
    this.eventsService.LoadEvents(50).subscribe((res) => {
      if (res) {
        this.dataSource = new MatTableDataSource(res);
        this.dataSource.sort = this.sort;
        this.loading = false;
      }
    });
  }
}
