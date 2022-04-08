import { Component, OnInit, Output } from '@angular/core';
import { EventEmitter } from '@angular/core';
import { ToolbarService } from './toolbar.service';

@Component({
  selector: 'app-toolbar',
  templateUrl: './toolbar.component.html',
  styleUrls: ['./toolbar.component.scss'],
})
export class ToolbarComponent implements OnInit {
  @Output() sidenavToggle = new EventEmitter();
  title: string;
  constructor(private toolbarService: ToolbarService) {}

  ngOnInit(): void {
    this.title = this.toolbarService.getTitle();
  }
  toggleNav() {
    this.sidenavToggle.emit();
  }
}
