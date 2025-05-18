---
title: journal-007
date: 2025-05-16
sn_id: 982124
---

Dear journal,

something has changed. I could have easily skipped you today, but I didn't.

I just came back from walking for an hour toward Austin's skyline from the remainder of a party. Since it's cockroach weather, I got pretty sweaty, but that was fine. It was part of the deal of getting a nice long walk in for myself, and I was looking forward to the shower. I also almost took a nice photo of my shadow, but it wasn't nice enough to save. Still working on my photography skills.

Before I went to the end of the party, I went to watch _Final Destination Bloodlines_ with [@Car](https://stacker.news/Car). It was a very good movie. I was worried I’d be too scared—sometimes I get spooked by my own hair—but it was juuust right. I even clapped at the end to see if anybody else would clap. Nobody else clapped, but that was fine, because I clapped for all of them.

---

Today, I fixed a very annoying bug with our Nostr embeds ([#2151](https://github.com/stackernews/stacker.news/issues/2151)). When you clicked on 'show full note,' it would only work for some notes. The others wouldn’t expand to show the full note but instead just displayed a scrollbar. The embed stayed so tiny, the text was literally unreadable.

When I looked into the issue again with a fresh perspective, I noticed a key difference between the Nostr embeds where the button worked and those where it didn’t: the broken ones contained a redirect. This broke the embed because the query parameter `embed=true` would get lost during the redirect. I submitted a [PR](https://github.com/fiatjaf/njump/pull/107) to fix this in [njump](https://njump.me/) and [@fiatjaf](https://stacker.news/fiatjaf) merged and deployed it right away. The beauty of FOSS.

```diff
  // if we originally got a note code or an nevent with no hints
  // augment the URL to point to an nevent with hints -- redirect
  if p, ok := decoded.(nostr.EventPointer); (ok && p.Author == "" && len(p.Relays) == 0) || prefix == "note" {
-   http.Redirect(w, r, "/"+data.nevent, http.StatusFound)
+   url := "/" + data.nevent
+   if r.URL.RawQuery != "" {
+     url += "?" + r.URL.RawQuery
+   }
+   http.Redirect(w, r, url, http.StatusFound)
    return
  }
```

I also fixed the following migration error that I mentioned [yesterday](/journal/006):

```
app  | 2025-05-16T20:56:53.127625000Z Applying migration `20250516042253_wallet_v2`
app  | 2025-05-16T20:56:53.202274000Z Error: P3018
app  | 2025-05-16T20:56:53.202294000Z
app  | 2025-05-16T20:56:53.202309000Z A migration failed to apply. New migrations cannot be applied before the error is recovered from. Read more about how to resolve migration issues in a production database: https://pris.ly/d/migrate-resolve
app  | 2025-05-16T20:56:53.202322000Z
app  | 2025-05-16T20:56:53.202337000Z Migration name: 20250516042253_wallet_v2
app  | 2025-05-16T20:56:53.202349000Z
app  | 2025-05-16T20:56:53.202365000Z Database error code: 23503
app  | 2025-05-16T20:56:53.202379000Z
app  | 2025-05-16T20:56:53.202392000Z Database error:
app  | 2025-05-16T20:56:53.202406000Z ERROR: insert or update on table "Withdrawl" violates foreign key constraint "Withdrawl_walletId_fkey"
app  | 2025-05-16T20:56:53.202427000Z DETAIL: Key (walletId)=(29) is not present in table "ProtocolWallet".
app  | 2025-05-16T20:56:53.202443000Z
app  | 2025-05-16T20:56:53.202458000Z DbError { severity: "ERROR", parsed_severity: Some(Error), code: SqlState(E23503), message: "insert or update on table \"Withdrawl\" violates foreign key constraint \"Withdrawl_walletId_fkey\"", detail: Some("Key (walletId)=(29) is not present in table \"ProtocolWallet\"."), hint: None, position: None, w
here_: None, schema: Some("public"), table: Some("Withdrawl"), column: None, datatype: None, constraint: Some("Withdrawl_walletId_fkey"), file: Some("ri_triggers.c"), line: Some(2608), routine: Some("ri_ReportViolation") }
app  | 2025-05-16T20:56:53.202498000Z
```

The migration starts working again when I run the migrations in a different order. During a rebase, the migration to delete duplicate wallets was between my changes to the vault and wallet schema. So all I had to do to fix this was change the name of the vault migration such that it comes after the migration to delete duplicate wallets.

I was unable to find out why that is the case.