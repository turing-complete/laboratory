function approximate
  use('Interaction');

  filename = locate('approximate');
  solution = h5read(filename, '/solution');
  surrogate = solution.Surrogate;

  no = double(surrogate.Outputs);
  surpluses = reshape(surrogate.Surpluses{1}, no, []);

  Plot.figure(1200, 400);
  for i = 1:no
    subplot(no, 1, i);
    semilogy(abs(surpluses(i, :)));
	Plot.title('Output %d', i-1);
  end
end
