with Text_Io;
use Text_Io;

package body productors_consumidors is

   protected body Monitor_Productor_Consumidor is

      procedure iniciar_maquina(n_clients: Integer) is
      begin

         max_refrescs_maquina := 10;
         n_refrescs_maquina := 0;
         clients_acabats := 0;
         total_clients := n_clients;

         Put_Line("*** La maquina esta preparada");

      end iniciar_maquina;

      entry Consumir (n_refresc_consumit: Integer; n_refrescs_a_consumir: Integer; nom_client: String) when n_refrescs_maquina > 0 is
      begin

         n_refrescs_maquina := n_refrescs_maquina - 1;
         Put_Line("--- " & nom_client & " agafa el refresc " &Integer'Image(n_refresc_consumit)& " /" &Integer'Image(n_refrescs_a_consumir)& ", li queden a la maquina" &Integer'Image(n_refrescs_maquina));

      end Consumir;

      entry Reposar(n_reposador: Integer; fi_reposador: out boolean) when n_refrescs_maquina /= max_refrescs_maquina or total_clients = clients_acabats is

         n_refrescs_a_reposar: integer := max_refrescs_maquina - n_refrescs_maquina;

      begin

         if(total_clients = clients_acabats) then

            Put_Line("+++ El reposador" &Integer'Image(n_reposador)& " diu: No hi ha clients, me'n vaig");
            fi_reposador := True;

         else

            n_refrescs_maquina := max_refrescs_maquina;
            Put_Line("+++ El reposador" &Integer'Image(n_reposador)& " reposa" &Integer'Image(n_refrescs_a_reposar)& " refrescs, ara hi ha" &Integer'Image(n_refrescs_maquina));
            fi_reposador := False;

         end if;

      end Reposar;

      entry fi_client(clients_restants: out integer) when total_clients /= clients_acabats is
      begin

         clients_acabats := clients_acabats + 1;
         clients_restants := total_clients - clients_acabats;

      end fi_client;

   end Monitor_Productor_Consumidor;

end productors_consumidors;
