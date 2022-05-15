import { Component, OnInit, Output } from "@angular/core";
import { EventEmitter } from "@angular/core";
import { ToolbarService } from "./toolbar.service";

@Component({
  selector: "app-toolbar",
  templateUrl: "./toolbar.component.html",
  styleUrls: ["./toolbar.component.scss"],
})
export class ToolbarComponent implements OnInit {
  @Output() sidenavToggle = new EventEmitter();
  title: string;
  envButtons = [
    { url: "http://192.168.10.159:8080", descr: "PROD", icon: "dns" },
    { url: "http://192.168.10.159:8081", descr: "QA", icon: "dns" },
    { url: "http://localhost:4200", descr: "LOCAL", icon: "dns" },
  ];
  constructor(private toolbarService: ToolbarService) {}

  ngOnInit(): void {
    this.title = this.toolbarService.getTitle();
  }
  toggleNav() {
    this.sidenavToggle.emit();
  }
}
