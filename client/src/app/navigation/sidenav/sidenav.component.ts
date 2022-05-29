import { Component, OnInit } from "@angular/core";

@Component({
  selector: "app-sidenav",
  templateUrl: "./sidenav.component.html",
  styleUrls: ["./sidenav.component.scss"],
})
export class SidenavComponent implements OnInit {
  constructor() {}
  menuItems = [
    { url: "/dashboard", descr: "Dashboard", icon: "dashboard" },
    { url: "/devices", descr: "Devices", icon: "sensors" },
    { url: "/datasets", descr: "Datasets", icon: "assignment" },
    { url: "/controllers", descr: "Controllers", icon: "extension" },
    { url: "/events", descr: "System events", icon: "speaker_notes" },
    { url: "/system", descr: "System information", icon: "info" },
  ];
  ngOnInit(): void {}
}
