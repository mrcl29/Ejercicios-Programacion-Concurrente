with Text_Io; use Text_Io;

package body maquinarefrescs is

    protected body Monitor_Productor_Consumidor is

        procedure iniciar (n_clients : Integer) is
        begin
            -- Inicialitzacio de les variables del monitor
            max_refrescs_maquina := 10;
            n_refrescs_maquina := 0;
            clients_acabats := 0;
            total_clients := n_clients;

            Put_Line ("********** La maquina esta preparada");

        end iniciar;

        entry Reposar (id : Integer; fi_reposador : out boolean)
          when n_refrescs_maquina /= max_refrescs_maquina
          or total_clients = clients_acabats
        is
            -- Calcul del nombre de refrescs a reposar
            n_refrescs_a_reposar : integer :=
              max_refrescs_maquina - n_refrescs_maquina;

        begin
            -- Comprovem si tots els clients han acabat
            if (total_clients = clients_acabats) then
                -- Si no hi ha mes clients, el reposador acaba
                Put_Line
                  ("++++++++++ El reposador"
                   & Integer'Image (id)
                   & " diu: No hi ha clients me'n vaig");
                fi_reposador := True;

            else
                -- Reposem els refrescs i actualitzem el comptador
                n_refrescs_maquina := max_refrescs_maquina;
                Put_Line
                  ("++++++++++ El reposador"
                   & Integer'Image (id)
                   & " reposa"
                   & Integer'Image (n_refrescs_a_reposar)
                   & " refrescs, ara n'hi ha"
                   & Integer'Image (n_refrescs_maquina));
                fi_reposador := False;

            end if;

        end Reposar;

        entry Agafar
          (n_refresc_agafats     : Integer;
           n_refrescs_a_consumir : Integer;
           nom_client            : String)
          when n_refrescs_maquina > 0
        is
        begin
            -- Reduim el nombre de refrescs a la maquina
            n_refrescs_maquina := n_refrescs_maquina - 1;
            Put_Line
              ("---------- "
               & nom_client
               & " agafa el refresc"
               & Integer'Image (n_refresc_agafats)
               & " /"
               & Integer'Image (n_refrescs_a_consumir)
               & " a la maquina en queden"
               & Integer'Image (n_refrescs_maquina));

        end Agafar;

        entry fi_client (clients_restants : out integer)
          when total_clients /= clients_acabats
        is
        begin
            -- Incrementem el comptador de clients acabats
            clients_acabats := clients_acabats + 1;
            -- Calculem els clients que encara no han acabat
            clients_restants := total_clients - clients_acabats;

        end fi_client;

    end Monitor_Productor_Consumidor;

end maquinarefrescs;
