---
title: wtf, nix is magic
date: 2024-04-20
sn_id: 512757
tags: nixos
---

I just installed NixOS and entered this into my /etc/nixos/configuration.nix:

```
services.xserver.enable = true;
services.xserver.windowManager.dwm.enable = true;
```

and then I ran `nixos-rebuild switch` followed by `systemctl start display-manager` and now I have `dwm` running?

No need to manually install `xorg-server`, `dwm` and friends, fiddling to run `startx` on boot while making sure keyboard mappings are correct?[^1]

wtf, nix is magic

[^1]: `services.xserver.xkb.layout` got you covered
